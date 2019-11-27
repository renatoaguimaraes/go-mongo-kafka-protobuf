package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/renatoaguimaraes/go-mongo-kafka-protobuf/kafka"
	"github.com/renatoaguimaraes/go-mongo-kafka-protobuf/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func monitor(done chan bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	done := make(chan bool, 1)
	monitor(done)

	// protobuffer model
	test := &model.Test{
		Label: proto.String("hello"),
		Type:  proto.Int32(17),
		Reps:  []int64{1, 2, 3},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	// publishing messages
	producer := make(chan []byte)
	defer close(producer)
	go kafka.Producer("test", producer)
	go func() {
		for i := 1; i <= 100; i++ {
			producer <- data
			// time.Sleep(2 * time.Second)
		}
	}()

	// cunsuming messages
	messages := make(chan []byte)
	defer close(messages)
	go kafka.Consumer("test", messages)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		defer cancel()
		if err != nil {
			log.Fatal("Error to connect mongodb")
		}
		collection := client.Database("testing").Collection("numbers")

		for message := range messages {
			newTest := &model.Test{}
			err = proto.Unmarshal(message, newTest)
			if err != nil {
				log.Fatal("unmarshaling error: ", err)
			}
			_, err := collection.InsertOne(ctx, bson.M{"label": newTest.GetLabel()})
			if err != nil {
				log.Fatal("Error to insert data")
			}
		}
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}

// jobs := make(chan string, 64)
// 	defer close(jobs)

// 	results := make(chan string, 64)
// 	defer close(results)

// 	for i := 1; i <= 64; i++ {
// 		go worker(jobs, results)
// 	}
// 	go process(jobs)
// 	go handler(results)
