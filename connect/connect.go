package connect

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DB holds our mongo client access for Client,Collection,Database
type DB struct {
	Col *mongo.Collection
}

//Cluster connection to our mongo cluster
func Cluster(url string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

//DataBase ...
func DataBase(c *mongo.Client, databaseName, collectionName string) *DB {
	return &DB{
		Col: c.Database(databaseName).Collection(collectionName),
	}
}

//FindbyID ...
func (db *DB) FindbyID(id string, data interface{}) (interface{}, error) {

	idNew, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = db.Col.FindOne(context.TODO(), bson.M{"_id": idNew}).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//FindAny ...
func (db *DB) FindAny(field, value string, data interface{}) (interface{}, error) {

	err := db.Col.FindOne(context.TODO(), bson.M{field: value}).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
