package greet_client

import (
	"applinh/gogrpcudemy/greet/greetpb"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func StartClient() {

	opts := grpc.WithInsecure()
	tls := false
	if tls {
		certFile := "ssl/ca.crt"
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "localhost")
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to connect %v \n", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDrectStreaming(c)
	//doUnaryDeadline(c, 5*time.Second) // will complete
	//doUnaryDeadline(c, 1*time.Second) // will fail

}

//------------------------------------------
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
	fmt.Printf("Long greet response: %v \n", res.Result)

}
func doBiDrectStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting BiDi streaming...")

	// we create a stream by invoking a client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream %v \n", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
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

	waitChan := make(chan struct{})

	// send messages to the server
	go func() {
		for _, v := range requests {
			fmt.Printf("Sending message: %v \n", v)
			stream.Send(v)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	// receive messages from the server
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error occured while reading server stream %v \n", err)
				break
			}
			fmt.Printf("Received: %v \n", msg.Result)
		}
		close(waitChan)
	}()

	// block until everything is done
	<-waitChan

}
func doUnaryDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting unary with deadline...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Thomas",
			LastName:  "Martin",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {

		if statusErr, ok := status.FromError(err); ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Println("Deadline exceeded.")
			} else {
				log.Printf("Unexepected error: %v", err)
			}
		} else {
			log.Fatalf("Error while calling GreetWithDeadlineRequest %v \n", err)
		}
		return
	}
	log.Printf("Response from GreetWithDeadlineRequest: %v \n", res.Result)

}
