package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//User ...
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserName string             `bson:"username"`
	Address  UserAddress        `bson:"address"`
	Phone    []string           `bson:"phone"`
}

//UserAddress ...
type UserAddress struct {
	PhysicalAddress string `bson:"physical_address"`
	Road            string `bson:"road"`
}

//Todo ...
type Todo struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID primitive.ObjectID `bson:"user_id"`
}

func connect() *mongo.Client {
	client, err := NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	return client
}

func main() {
	client := connect()
	//insert user
	u, err := addUser(client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user details:", u)

}

func addUser(c *mongo.Client) (*User, error) {
	u := User{
		UserName: "emmanuel",
	}
	col := c.Database("mongolang").Collection("users")
	res, err := col.InsertOne(context.TODO(), u)
	if err != nil {
		return nil, err
	}
	resUser := User{}
	id, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", res.InsertedID))

	_ = col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&resUser)

	return &resUser, nil
}
