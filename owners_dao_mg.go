package main

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OwnersDao interface {
	GetAll() ([]*Owner, error)
	Get(id string) (*Owner, error)
	Create(owner *Owner) error
	Update(owner *Owner) error
	Delete(id string) error
}

type OwnerDaoImpl struct {
	client *mongo.Client
}

func (dao *OwnerDaoImpl) GetAll() ([]*Owner, error) {
	ctx := context.Background()
	owners := []*Owner{}
	collection := dao.client.Database("todos").Collection("owners")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil && err.Error() == "document is nil" {
		log.Printf("could not get owners collection: %s", err.Error())
		return owners, nil
	} else if err != nil {
		return nil, errors.New("could not get owners collection: " + err.Error())
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &owners)
	if err != nil {
		return nil, errors.New("could not get owners: " + err.Error())
	}
	return owners, nil
}

func (dao *OwnerDaoImpl) Get(id string) (*Owner, error) {
	ctx := context.Background()
	owner := &Owner{}
	collection := dao.client.Database("todos").Collection("owners")
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(owner)
	if err != nil {
		return nil, errors.New("could not get owner: " + err.Error())
	}
	return owner, nil
}

func (dao *OwnerDaoImpl) Create(owner *Owner) error {
	ctx := context.Background()
	collection := dao.client.Database("todos").Collection("owners")
	_, err := collection.InsertOne(ctx, owner)
	if err != nil {
		return errors.New("could not create owner: " + err.Error())
	}
	return nil
}

func (dao *OwnerDaoImpl) Update(owner *Owner) error {
	ctx := context.Background()
	collection := dao.client.Database("todos").Collection("owners")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": owner.ID}, bson.M{"$set": owner})
	if err != nil {
		return errors.New("could not update owner: " + err.Error())
	}
	return nil
}

func (dao *OwnerDaoImpl) Delete(id string) error {
	ctx := context.Background()
	collection := dao.client.Database("todos").Collection("owners")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.New("could not delete owner: " + err.Error())
	}
	return nil
}

func NewOwnerDao(client *mongo.Client) OwnersDao {
	return &OwnerDaoImpl{client}
}