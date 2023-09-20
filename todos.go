package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gtihub.com/google/uuid"
)

type Todo struct {
	ID        int       `bson:"_id"`
	Task      string    `bson:"task"`
	Completed bool      `bson:"completed"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type TodosDao interface {
	GetAll() ([]*Todo, error)
	Get(id int) (*Todo, error)
	Create(todo *Todo) error
	Update(todo *Todo) error
	Delete(id int) error
}

type TodoDaoImpl struct {
	client *mongo.Client
}

func NewTodoDao(client *mongo.Client) TodosDao {
	return &TodoDaoImpl{client}
}

func (dao *TodoDaoImpl) GetAll() ([]*Todo, error) {
	ctx := context.Background()
	// get all todos from mongodb
	coll := dao.client.Database("todos").Collection("todos")
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// convert todos to []*Todo
	var todos []*Todo
	for cur.Next(ctx) {
		var todo *Todo
		if err := cur.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (dao *TodoDaoImpl) Get(id int) (*Todo, error) {
	ctx := context.Background()
	todo := &Todo{}
	// get todo from mongodb
	coll := dao.client.Database("todos").Collection("todos")
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(todo)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (dao *TodoDaoImpl) Create(todo *Todo) error {
	todo.ID = uuid.New().String()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	_, err := dao.client.Database("todos").Collection("todos").InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}
	return nil
}

func (dao *TodoDaoImpl) Update(todo *Todo) error {
	return errors.New("not implemented")
}

func (dao *TodoDaoImpl) Delete(id int) error {
	return errors.New("not implemented")
}
