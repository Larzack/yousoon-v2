module github.com/yousoon/services/partner

go 1.21

require (
	github.com/google/uuid v1.5.0
	github.com/yousoon/services/shared v0.0.0
	go.mongodb.org/mongo-driver v1.13.1
)

replace github.com/yousoon/services/shared => ../shared
