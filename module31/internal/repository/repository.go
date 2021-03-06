package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"skillbox/internal/entity"
)

var collection *mongo.Collection
var ctx = context.TODO()

type mongodb struct {
	client *mongo.Client
}

func NewMongodb() (*mongodb, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &mongodb{client: client}, nil
}

func DisconnectDB(client *mongodb) {
	err := client.client.Disconnect(ctx)
	if err != nil {
		panic(err)
	}
}

func (r *mongodb) CreateUser(user *entity.User) (string, error) {
	collection = r.client.Database("usersDB").Collection("users")
	u, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Println(err)
		return "", err
	}

	id := u.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (r *mongodb) DeleteUser(id string) (string, error) {
	collection = r.client.Database("usersDB").Collection("users")
	userID, err := primitive.ObjectIDFromHex(id)
	var d bson.M
	_ = collection.FindOneAndDelete(ctx, bson.D{{
		"_id",
		userID,
	}}).Decode(&d)

	name := d["Name"].(string)
	if err != nil {
		log.Println(err)
		return "", err
	}
	filter := bson.D{{"Friends", id}}
	fmt.Println(name)
	update := bson.D{
		{"$pull", bson.D{
			{"Friends", id},
		}},
	}
	_, err = collection.UpdateMany(ctx, filter, update)
	return name, nil
}

func (r *mongodb) GetUsers(user *entity.User) []*entity.User {
	collection = r.client.Database("usersDB").Collection("users")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		panic(err)
	}
	var allUsers []*entity.User
	for cur.Next(ctx) {

		err := cur.Decode(&user)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		allUsers = append(allUsers, user)
	}
	err = cur.Close(ctx)
	return allUsers
}

func (r *mongodb) UpdateAge(id string, newAge int) error {
	collection = r.client.Database("usersDB").Collection("users")
	userID, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", userID}}
	update := bson.D{
		{"$set", bson.D{
			{"Age", newAge},
		}},
	}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *mongodb) MakeFriends(target string, source string) (string, string, error) {
	collection = r.client.Database("usersDB").Collection("users")
	targetID, _ := primitive.ObjectIDFromHex(target)
	sourceID, _ := primitive.ObjectIDFromHex(source)
	opt := bson.D{
		{"_id", 0},
		{"Name", 1},
	}
	cur, _ := collection.Find(ctx, bson.D{{
		"_id",
		bson.D{{
			"$in",
			bson.A{targetID, sourceID},
		}},
	}}, options.Find().SetProjection(opt))
	var n bson.M
	var names []string

	for cur.Next(ctx) {
		_ = cur.Decode(&n)
		names = append(names, n["Name"].(string))
	}

	filter := bson.D{{"_id", targetID}}
	update := bson.D{
		{"$push", bson.D{
			{"Friends", source},
		}},
	}
	_, _ = collection.UpdateOne(ctx, filter, update)
	filter = bson.D{{"_id", sourceID}}
	update = bson.D{
		{"$push", bson.D{
			{"Friends", target},
		}},
	}
	_, _ = collection.UpdateOne(ctx, filter, update)
	return names[0], names[1], nil
}

func (r *mongodb) GetFriends(userId string) ([]string, error) {
	collection = r.client.Database("usersDB").Collection("users")
	var user bson.M

	f, err := collection.Find(ctx, bson.D{{"Friends", userId}})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var friends []string
	for f.Next(ctx) {
		_ = f.Decode(&user)
		friends = append(friends, user["Name"].(string))
	}
	return friends, nil
}
