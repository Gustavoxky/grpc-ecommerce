FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o grpc-server ./main.go

EXPOSE 50051

CMD ["./grpc-server"]