package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/golang/protobuf/proto"
	"github.com/renatoaguimaraes/go-mongo-kafka-protobuf/kafka"
	"github.com/renatoaguimaraes/go-mongo-kafka-protobuf/model"
	"github.com/renatoaguimaraes/go-mongo-kafka-protobuf/repo"
	"github.com/renatoaguimaraes/go-mongo-kafka-protobuf/util"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	done := make(chan bool, 1)
	util.Monitor(done)

	// protobuffer model
	test := &model.Test{Label: proto.String("hello"), Type: proto.Int32(17), Reps: []int64{1, 2, 3}}
	testProto, err := proto.Marshal(test)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	// publishing messages
	producer := make(chan []byte, 20)
	defer close(producer)
	go kafka.Producer("test", producer)
	go func() {
		for i := 1; i <= 10000; i++ {
			producer <- testProto
		}
	}()

	messages := make(chan []byte, 20)
	defer close(messages)

	for i := 1; i <= 20; i++ {
		go kafka.Consumer("test", messages)
	}

	for i := 1; i <= 20; i++ {
		go repo.InsertTest(messages)
	}

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
