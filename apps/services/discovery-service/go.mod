module github.com/yousoon/discovery-service

go 1.23

require (
	github.com/99designs/gqlgen v0.17.42
	github.com/elastic/go-elasticsearch/v8 v8.11.1
	github.com/go-chi/chi/v5 v5.0.11
	github.com/google/uuid v1.5.0
	github.com/gorilla/websocket v1.5.1
	github.com/nats-io/nats.go v1.31.0
	github.com/prometheus/client_golang v1.18.0
	github.com/redis/go-redis/v9 v9.3.1
	github.com/vektah/gqlparser/v2 v2.5.10
	github.com/yousoon/shared v0.0.0
	go.mongodb.org/mongo-driver v1.13.1
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.opentelemetry.io/otel/sdk v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
	go.uber.org/zap v1.26.0
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sosodev/duration v1.1.0 // indirect
	github.com/urfave/cli/v2 v2.25.5 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.9.3 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/yousoon/shared => ../shared
