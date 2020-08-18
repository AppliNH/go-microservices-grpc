FROM golang:latest

ENV GO111MODULE=on

WORKDIR /app

COPY ./go.mod .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o grpc_services .

EXPOSE 50051 50052 50053

CMD ["./grpc_services"]