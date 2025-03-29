run-server:
	go run main.go

run-client:
	go run client/main.go

proto:
	protoc --go_out=. --go-grpc_out=. proto/order.proto

tidy:
	go mod tidy

build:
	go build -o grpc-server main.go

docker-build:
	docker build -t grpc-ecommerce-server .

docker-run:
	docker run -p 50051:50051 grpc-ecommerce-server
