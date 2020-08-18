package greet_client

import (
	"applinh/gogrpcudemy/greet/greetpb"
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	doClientStreaming(c)

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
	fmt.Println("Starting server streaming...")

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

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting client streaming...")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet %v", err)
	}

	requests := []*greetpb.LongGreetRequest{
		{ //typing &greetpb.LongGreetRequest is redudant here and not needed
			Greeting: &greetpb.Greeting{
				FirstName: "Thomas",
				LastName:  "Martin",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Richard",
				LastName:  "Dupont",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Brett",
				LastName:  "Fisher",
			},
		},
	}

	for _, v := range requests {
		fmt.Println("Sending for: " + v.Greeting.FirstName)
		stream.Send(v)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving from server %v", err)
	}
	fmt.Printf("Long greet response: %v", res.Result)

}
