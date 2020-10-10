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

//FindbyObjectID ...
//find by object id string
func (db *DB) FindbyObjectID(idObjectString string) (bson.M, error) {
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

//FindAll ...
//with multiple filters using primitive.M
func (db *DB) FindAll(filter bson.M) ([]bson.M, error) {
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

	allResults, err := db.FindAll(bson.M{"_id": bson.M{"$in": IDs}})
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

	result, err := db.FindOne(bson.M{"_id": res.InsertedID})
	if err != nil {
		return nil, err
	}

	return result, nil
}

//CreateM ...
func (db *DB) CreateM(insertData bson.M) (bson.M, error) {

	res, err := db.Col.InsertOne(context.TODO(), insertData)
	if err != nil {
		return nil, err
	}
	result, err := db.FindOne(bson.M{"_id": res.InsertedID})
	if err != nil {
		return nil, err
	}

	return result, nil
}

//Update ...
func (db *DB) Update(filter bson.M, updateParams bson.M) (bson.M, error) {
	res, err := db.Col.UpdateOne(context.Background(), filter, updateParams)
	if err != nil {
		return nil, err
	}

	result, err := db.FindOne(bson.M{"_id": res.UpsertedID})
	if err != nil {
		return nil, err
	}

	return result, nil
}

//UpdateByID ...
func (db *DB) UpdateByID(IDObject string, updateParams bson.M) (bson.M, error) {
	idNew, _ := primitive.ObjectIDFromHex(IDObject)
	filter := bson.M{
		"_id": idNew,
	}
	res, err := db.Col.UpdateOne(context.Background(), filter, updateParams)
	if err != nil {
		return nil, err
	}

	result, err := db.FindOne(bson.M{"_id": res.UpsertedID})
	if err != nil {
		return nil, err
	}

	return result, nil
}

//UpdateMany ...
func (db *DB) UpdateMany(filter bson.M, updateParams bson.M) (int64, error) {
	res, err := db.Col.UpdateMany(context.Background(), filter, updateParams)
	if err != nil {
		return 0, err
	}

	return res.UpsertedCount, nil

}

//SoftDelete ...
func (db *DB) SoftDelete(filter bson.M) (bson.M, error) {
	deletedAt := time.Now()
	res, err := db.Col.UpdateOne(context.Background(), filter, bson.M{"$set": bson.M{"deleted_at": deletedAt}})
	if err != nil {
		return nil, err
	}
	result, err := db.FindOne(bson.M{"_id": res.UpsertedID})
	if err != nil {
		return nil, err
	}

	return result, nil
}

//Delete ...
func (db *DB) Delete(filter bson.M) (int64, error) {
	res, err := db.Col.DeleteOne(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

//DeleteMany ..
func (db *DB) DeleteMany(filter bson.M) (int64, error) {
	res, err := db.Col.DeleteMany(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
