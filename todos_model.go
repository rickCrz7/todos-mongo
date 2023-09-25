package main

import "time"

type Owner struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type Todo struct {
	ID        string    `json:"id" bson:"_id"`
	Title     string    `json:"title" bson:"title"`
	Completed bool      `json:"completed" bson:"completed"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	OwnerID   string    `json:"owner_id" bson:"owner_id"`
	Owner     *Owner    `json:"owner" bson:"owner"`
}