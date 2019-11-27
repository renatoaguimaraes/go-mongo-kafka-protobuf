```shell
go get "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
go get "github.com/golang/protobuf/proto"
go get "go.mongodb.org/mongo-driver/mongo"
```
```shell
protoc --go_out=. model/*.proto
```
```shell
docker-compose -f docker-compose-dev.yml up
```
