module nexa/services/authorization

go 1.22.4

require (
	github.com/caarlos0/env/v10 v10.0.0
	github.com/uptrace/bun v1.2.1
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.2
	nexa/proto/gen/go v0.0.0
	nexa/shared v0.0.0
)

replace (
	nexa/proto/gen/go => ../../proto/gen/go
	nexa/shared => ../../shared/
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240617180043-68d350f18fd4 // indirect
)
