package connect

import (
	"context"
	"fmt"
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

//FindOnebyID ...
//find by object id string
func (db *DB) FindOnebyID(idObjectString string) (bson.M, error) {
	var data bson.M
	idNew, err := primitive.ObjectIDFromHex(idObjectString)
	if err != nil {
		return nil, err
	}

	err = db.Col.FindOne(context.TODO(), bson.M{"_id": idNew}).Decode(&data)
	if err != nil {
		fmt.Println("found nothing here")
		return nil, err
	}

	return data, nil
}

//FindOneByField ...
//user any field but not object ids
func (db *DB) FindOneByField(field, value string) (bson.M, error) {
	var data bson.M
	err := db.Col.FindOne(context.TODO(), bson.M{field: value}).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//FindOne ...
//with multiple filters using primitive.M
func (db *DB) FindOne(filter bson.M) (bson.M, error) {
	var data bson.M
	err := db.Col.FindOne(context.TODO(), filter).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//FindMany ...
//with multiple filters using primitive.M
func (db *DB) FindMany(filter bson.M) ([]bson.M, error) {
	var ManyResults []bson.M
	cursor, err := db.Col.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		singleReslut := bson.M{}
		err := cursor.Decode(&singleReslut)
		if err != nil {
			return nil, err
		}
		ManyResults = append(ManyResults, singleReslut)
	}

	return ManyResults, nil
}

//CreateMany ...
func (db *DB) CreateMany(insertData []interface{}) ([]bson.M, error) {
	var order bool = true
	res, err := db.Col.InsertMany(context.Background(), insertData, &options.InsertManyOptions{Ordered: &order})
	if err != nil {
		return nil, err
	}

	var IDs []interface{}
	for _, id := range res.InsertedIDs {
		IDs = append(IDs, id)
	}

	allResults, err := db.FindMany(bson.M{"_id": bson.M{"$in": IDs}})
	if err != nil {
		return nil, err
	}

	return allResults, nil
}

//Create ...
func (db *DB) Create(insertData interface{}) (bson.M, error) {

	res, err := db.Col.InsertOne(context.TODO(), insertData)
	if err != nil {
		return nil, err
	}

	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", res.InsertedID))
	if err != nil {
		return nil, err
	}

	results, err := db.FindOne(bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	return results, nil
}
