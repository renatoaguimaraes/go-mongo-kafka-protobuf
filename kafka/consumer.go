package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

//Consumer kafka messages
func Consumer(topic string, messages chan []byte) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	defer c.Close()

	if err != nil {
		panic(err)
	}
	c.SubscribeTopics([]string{topic}, nil)
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			messages <- msg.Value
		} else {
			// TODO handle errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
