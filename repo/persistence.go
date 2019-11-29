package repo

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/renatoaguimaraes/go-mongo-kafka-protobuf/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertTest insert a
func InsertTest(messages chan []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer cancel()
	if err != nil {
		log.Fatal("Error to connect mongodb")
	}
	for message := range messages {
		collection := client.Database("testing").Collection("numbers")
		test := &model.Test{}
		err = proto.Unmarshal(message, test)
		if err != nil {
			log.Fatal("Unmarshaling error: ", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := collection.InsertOne(ctx, test)
		if err != nil {
			log.Fatal("Error to insert document: ", err)
		}
		cancel()
	}
}
