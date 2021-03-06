package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/techcodemy/mongolang/connect"
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
	Phone    []string           `bson:"phone,omitempty"`
	Todos    []Todo             `bson:"todos,omitempty"`
}

//UserAddress ...
type UserAddress struct {
	PhysicalAddress string `bson:"physical_address"`
	Road            string `bson:"road"`
}

//Todo ...
type Todo struct {
	UserID   primitive.ObjectID `bson:"user_id"`
	Text     string             `bson:"text"`
	Done     bool               `bson:"done"`
	Comments []Comment          `bson:"comments"`
}

//Comment ...
type Comment struct {
	UserID primitive.ObjectID `bson:"user_id"`
	Cotent string             `bson:"content"`
}

/* func connect() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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
} */

const url = "mongodb://localhost:27017"

func main() {
	/* connecting try connection */
	//client, _ := connect()

	client, err := connect.Cluster(url)
	if err != nil {
		log.Fatal(err)
	}
	db := connect.DataBase(client, "mongolang", "users")

	/* // #findByID
	res, err := db.FindOnebyID("5f7ef973ef75227b2f76e007")
	if err != nil {
		log.Fatal(err)
	}*/
	// #findOneByField
	/* res, err := db.FindOneByField("username", "esnart")
	if err != nil {
		log.Fatal(err)
	} */

	// #findOne

	/* res, err := db.FindOne(bson.M{"_id": id})
	if err != nil {
		log.Fatal(err)
	}

	user := User{}
	bsonBytes, _ := bson.Marshal(res)
	bson.Unmarshal(bsonBytes, &user)
	fmt.Println(user) */

	// getting all
	res, err := db.FindAll(bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	usersStruct := []User{}

	for _, i := range res {
		user := User{}
		bsonBytes, _ := bson.Marshal(i)
		bson.Unmarshal(bsonBytes, &user)

		usersStruct = append(usersStruct, user)
	}

	fmt.Println("return user:", usersStruct)

	//create using orm
	/* u := bson.M{
		"username": "emmunic boss",
	} */

	/* u := User{
		ID:       primitive.NewObjectID(),
		UserName: "user test",
	}
	res, err := db.Create(u)
	if err != nil {
		log.Fatal(err)
	}

	user := User{}
	bsonBytes, _ := bson.Marshal(res)
	bson.Unmarshal(bsonBytes, &user)
	fmt.Println(user) */

	/* mary := User{
		ID:       primitive.NewObjectID(),
		UserName: "mary",
		Phone:    &[]string{"0945934878"},
	}
	bobo := User{ID: primitive.NewObjectID(), UserName: "bobo"}

	res, err := db.CreateMany([]interface{}{mary, bobo})
	if err != nil {
		log.Fatal(err)
	}

	usersStruct := []User{}

	for _, i := range res {
		user := User{}
		bsonBytes, _ := bson.Marshal(i)
		bson.Unmarshal(bsonBytes, &user)

		usersStruct = append(usersStruct, user)
	}

	fmt.Println("return user:", usersStruct) */

	//Finding one
	/* uME := User{}
	idNew, _ := primitive.ObjectIDFromHex("5f7ef973ef75227b2f76e007")
	err = db.Col.FindOne(context.TODO(), bson.M{"_id": idNew}).Decode(&uME)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user struct:", uME) */

	//fmt.Println("user interface:", u)

	// ***insert user
	/* u, err := addUser(client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user details:", u) */

	//inserting multiple users
	/* us, err := addManyUser(client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("users lists:\n", us) */

	// updating or pushing to an array
	/* u, err := addPhoneToUser("5f7ef973ef75227b2f76e007", client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user updated:\n", u) */

	// updating embedded details
	/* u, err := addUserAddress("5f7ef973ef75227b2f76e007", client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user address updated:\n", u) */

	//adding a todo in an array
	/* todos, err := addNewTodo("5f7ef973ef75227b2f76e007", client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("added a todo:\n", todos) */

	//removing a todo from an array by its index
	/* todos, err := deleteATodoFromArray("5f7ef973ef75227b2f76e007", 2, client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("deleted a todo:\n", todos) */

	//updating make as read in array by its index
	/* todos, err := MarkAsDoneTodo("5f7ef973ef75227b2f76e007", 2, client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("deleted a todo:\n", todos) */

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

func addManyUser(c *mongo.Client) (*[]User, error) {
	col := c.Database("mongolang").Collection("users")
	naile := User{
		ID:       primitive.NewObjectID(),
		UserName: "naile",
		Phone:    []string{"0945934878"},
	}
	kamuel := User{ID: primitive.NewObjectID(), UserName: "kamuel"}
	nickolas := User{ID: primitive.NewObjectID(), UserName: "nickolas"}
	esnart := User{ID: primitive.NewObjectID(), UserName: "esnart"}
	samuel := User{ID: primitive.NewObjectID(), UserName: "samuel"}

	var order bool = true
	res, err := col.InsertMany(context.Background(), []interface{}{naile, kamuel, nickolas, esnart, samuel}, &options.InsertManyOptions{Ordered: &order})
	if err != nil {
		return nil, err
	}
	resUsers := []User{}
	var IDs []interface{}
	for _, id := range res.InsertedIDs {

		fmt.Println("check insert id:", id)
		/* localID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", id))
		fmt.Println("loc id:", localID) */
		IDs = append(IDs, id)

	}

	cursor, _ := col.Find(context.Background(), bson.M{"_id": bson.M{"$in": IDs}})
	for cursor.Next(context.Background()) {
		singleReslut := User{}
		err := cursor.Decode(&singleReslut)
		if err != nil {
			return nil, err
		}
		resUsers = append(resUsers, singleReslut)
	}

	return &resUsers, nil
}

func addPhoneToUser(id string, c *mongo.Client) (*User, error) {
	col := c.Database("mongolang").Collection("users")
	idNew, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": idNew,
	}

	// ***** updating if one value
	/* res, err := col.UpdateOne(context.Background(), filter,
		bson.M{"$set": bson.M{"phone": "0955404864"}},
	) */

	// ***** pushing on if an array (method 1 without struct)
	_, err := col.UpdateOne(context.Background(), filter,
		bson.M{"$push": bson.M{"phone": []interface{}{"0955404864", "0973827172"}}},
	)

	/* p := &[]string{
		"0975503928",
		"0975593672",
	} */
	//_, err := col.UpdateOne(context.Background(), filter, bson.M{"$push": bson.M{"phone": bson.A{p}}})

	if err != nil {
		log.Fatal(err)
	}

	resUser := User{}

	idNew2, _ := primitive.ObjectIDFromHex(id)

	_ = col.FindOne(context.TODO(), bson.M{"_id": idNew2}).Decode(&resUser)

	return &resUser, nil
}

func addUserAddress(id string, c *mongo.Client) (*UserAddress, error) {
	col := c.Database("mongolang").Collection("users")
	idNew, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": idNew,
	}

	change := UserAddress{
		PhysicalAddress: "chilenje south",
		Road:            "kalomo rd",
	}

	// ***** updating if one value
	_, err := col.UpdateOne(context.Background(), filter,
		bson.M{"$set": bson.M{"address": change}},
	)
	if err != nil {
		log.Fatal(err)
	}

	resUser := User{}

	_ = col.FindOne(context.TODO(), bson.M{"_id": idNew}).Decode(&resUser)

	return &resUser.Address, nil
}

func addNewTodo(userID string, c *mongo.Client) ([]Todo, error) {
	col := c.Database("mongolang").Collection("users")
	idNew, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{
		"_id": idNew,
	}

	todo := Todo{
		UserID: idNew,
		Text:   "talk to madam about the cream",
		Done:   false,
	}

	_, err := col.UpdateOne(context.TODO(), filter, bson.M{"$push": bson.M{"todos": todo}})

	if err != nil {
		log.Fatal(err)
	}

	resUser := User{}

	err = col.FindOne(context.TODO(), filter).Decode(&resUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("count array:", len(resUser.Todos))

	return resUser.Todos, nil
}

func deleteATodoFromArray(userID string, index int, c *mongo.Client) ([]Todo, error) {
	col := c.Database("mongolang").Collection("users")
	idNew, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{
		"_id": idNew,
	}

	fieldIndex := "todos." + strconv.Itoa(index)

	//unset removes everything in the array including the field
	//_, err := col.UpdateOne(context.TODO(), filter, bson.M{"$unset": bson.M{"todos": index}})

	_, err := col.UpdateOne(context.TODO(), filter, bson.M{"$unset": bson.M{fieldIndex: 0}})
	if err != nil {
		log.Fatal(err)
	}
	//removing all the null elecments in the database
	/* _, err = col.UpdateMany(context.TODO(), filter, bson.M{"$pull": bson.M{"todos": nil}})
	if err != nil {
		log.Fatal(err)
	} */
	resUser := User{}

	err = col.FindOne(context.TODO(), bson.M{"_id": idNew}).Decode(&resUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("count array:", len(resUser.Todos))

	return resUser.Todos, nil
}

func MarkAsDoneTodo(userID string, index int, c *mongo.Client) ([]Todo, error) {
	col := c.Database("mongolang").Collection("users")
	idNew, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{
		"_id": idNew,
	}
	findTodos := "todos." + strconv.Itoa(index) + ".done"

	_, err := col.UpdateOne(context.TODO(), filter,
		bson.M{"$set": bson.M{findTodos: true}})
	if err != nil {
		log.Fatal(err)
	}
	resUser := User{}

	err = col.FindOne(context.TODO(), bson.M{"_id": idNew}).Decode(&resUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("count array:", len(resUser.Todos))

	return resUser.Todos, nil
}
