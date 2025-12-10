module github.com/yousoon/apps/services/notification-service

go 1.21

require (
	github.com/99designs/gqlgen v0.17.45
	github.com/vektah/gqlparser/v2 v2.5.11
	go.mongodb.org/mongo-driver v1.14.0
	github.com/nats-io/nats.go v1.33.1
	github.com/aws/aws-sdk-go-v2 v1.25.2
	github.com/aws/aws-sdk-go-v2/service/ses v1.21.0
	github.com/aws/aws-sdk-go-v2/service/sns v1.29.0
	github.com/yousoon/shared v0.0.0
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/sosodev/duration v1.2.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.20.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/yousoon/shared => ../shared
