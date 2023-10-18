package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRepository struct {
	usersCollection *mongo.Collection
	chatroomsCollection *mongo.Collection
}

func NewMongoDBRepository(db *mongo.Database) *MongoDBRepository {
	return &MongoDBRepository{
		db.Collection("users"),
		db.Collection("chatrooms"),
	}
}

func (repo *MongoDBRepository) PopulateMongoDB() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	usersData := []bson.M{
		{
			"_id":         "652a9d83f896448755f4440f",
			"email":       "artichoke@jobsity.com",
			"password":    "$2a$14$ukT1jGx.NRm/giDivYzlN.N7jQP/U8ErNwxICu.FiqTXky5wm5GUS",
			"name":        "Artichoke",
			"currentRoom": "room2",
		},
		{
			"_id":         "652a9dd9f896448755f44410",
			"email":       "lilbasil@jobsity.com",
			"password":    "$2a$14$ukT1jGx.NRm/giDivYzlN.N7jQP/U8ErNwxICu.FiqTXky5wm5GUS",
			"name":        "Lil' Basil",
			"currentRoom": "room1",
		},
		{
			"_id":         "653013ed4cca2b12068bf08f",
			"email":       "StockBot",
			"password":    "",
			"name":        "StockBot",
			"currentRoom": nil,
		},
	}

	chatroomsData := []bson.M{
		{
			"_id":   "652d548a4cca2b12068bf08d",
			"name":  "Offtopic",
			"desc":  "Come here and make some friends",
			"code":  "room1",
		},
		{
			"_id":   "652d557c4cca2b12068bf08e",
			"name":  "Serious Talk",
			"desc":  "Discuss technical stuff",
			"code":  "room2",
		},
	}

	usersCollection := repo.usersCollection
	var usersDataInterface []interface{}
	for _, user := range usersData {
		usersDataInterface = append(usersDataInterface, user)
	}
	_, err = usersCollection.InsertMany(context.TODO(), usersDataInterface)
	if err != nil {
		return err
	}

	chatroomsCollection := repo.chatroomsCollection
	var chatroomsDataInterface []interface{}
	for _, chatroom := range chatroomsData {
		chatroomsDataInterface = append(chatroomsDataInterface, chatroom)
	}
	_, err = chatroomsCollection.InsertMany(context.TODO(), chatroomsDataInterface)
	if err != nil {
		return err
	}
	
	return nil
}