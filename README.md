## Demo golang, mongodb, kafka and protobuf

Installing dependencies.
```shell
go get "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
go get "github.com/golang/protobuf/proto"
go get "go.mongodb.org/mongo-driver/mongo"
```

Generating model files .go based on protobuf .
```shell
protoc --go_out=. model/*.proto
```

Starting kafka, zookeeper and mongodb.
```shell
docker-compose -f docker-compose-dev.yml up
```
