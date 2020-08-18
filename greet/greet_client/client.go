package greet_client

import (
	"applinh/gogrpcudemy/greet/greetpb"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func StartClient() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect %v \n", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doUnary(c)
	doServerStreaming(c)

}
func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting unary...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Thomas",
			LastName:  "Martin",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet %v \n", err)
	}
	log.Printf("Response from Greet: %v \n", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting streaming...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Thomas",
			LastName:  "Martin",
		},
	}
	res, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes %v \n", err)
	}

	for {
		msg, err := res.Recv()
		if err == io.EOF {
			// reached end of file
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream %v", err)
		}
		log.Printf("%v", msg.GetResult())
	}

}
