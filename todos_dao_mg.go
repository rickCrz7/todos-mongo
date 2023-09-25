package main

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodosDao interface {
	GetAll() ([]*Todo, error)
	Get(id string) (*Todo, error)
	Create(todo *Todo) error
	Update(todo *Todo) error
	Delete(id string) error
	Done(id string) error
}

type TodoDaoImpl struct {
	client *mongo.Client
}


func (dao *TodoDaoImpl) GetAll() ([]*Todo, error) {
	ctx := context.Background()
	todos := []*Todo{}
	collection := dao.client.Database("todos").Collection("todos")
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "owners",
				"localField":   "owner_id",
				"foreignField": "_id",
				"as":           "owner",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$owner",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$sort": bson.M{
				"created_at": -1,
			},
		},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

func (dao *TodoDaoImpl) Get(id string) (*Todo, error) {
	ctx := context.Background()
	todos := []*Todo{}
	collection := dao.client.Database("todos").Collection("todos")
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": id,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "owners",
				"localField":   "owner_id",
				"foreignField": "_id",
				"as":           "owner",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$owner",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &todos)
	if err != nil {
		return nil, err
	}
	if len(todos) == 0 {
		return nil, errors.New("todo " + id + " not found")
	}
	return todos[0], nil
}

func (dao *TodoDaoImpl) Create(todo *Todo) error {
	ctx := context.Background()
	insert := bson.M{
		"_id":        uuid.New().String(),
		"title":      todo.Title,
		"completed":  todo.Completed,
		"created_at": time.Now(),
		"updated_at": time.Now(),
		"owner_id":   todo.Owner.ID,
	}
	collection := dao.client.Database("todos").Collection("todos")
	if _, err := collection.InsertOne(ctx, insert); err != nil {
		return err
	}
	return nil
}

func (dao *TodoDaoImpl) Update(todo *Todo) error {
	ctx := context.Background()
	updatedAt := time.Now()
	collection := dao.client.Database("todos").Collection("todos")
	updateTodo := bson.M{
		"title":      todo.Title,
		"completed":  todo.Completed,
		"updated_at": updatedAt,
	}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": todo.ID}, bson.M{"$set": updateTodo})
	if err != nil {
		return err
	}
	return nil
}

func (dao *TodoDaoImpl) Delete(id string) error {
	ctx := context.Background()
	collection := dao.client.Database("todos").Collection("todos")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (dao *TodoDaoImpl) Done(id string) error {
	ctx := context.Background()
	collection := dao.client.Database("todos").Collection("todos")
	updatedAt := time.Now()
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"completed": true, "updated_at": updatedAt}})
	if err != nil {
		return err
	}
	return nil
}

func NewTodoDao(client *mongo.Client) TodosDao {
	return &TodoDaoImpl{client}
}