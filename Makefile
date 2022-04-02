run:
	clear
	go run cmd/grpc/main.go

tidy:
	go mod tidy

proto: build-proto inject-proto

build-proto:
	protoc --proto_path=./internal/core/domain:. --micro_out=. --go_out=. user.proto

inject-proto:
	protoc-go-inject-tag -input="./internal/core/domain/*.pb.go"

build-grpc:
	go build -o ./bin/user.grpc ./cmd/grpc/main.go

build-http:
	go build -o ./bin/user.http ./cmd/http/main.go

