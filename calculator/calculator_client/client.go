package calculator_client

import (
	"applinh/gogrpcudemy/calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func StartClient() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect %v \n", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDirectStreaming(c)
	doErrorUnary(c)

}
func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		A: 4,
		B: 5,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum %v \n", err)
	}
	log.Printf("Response from Sum: %v \n", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting streaming server")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	res, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition %v \n", err)
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
		log.Printf("%v", msg.GetPrimeNumber())
	}

}
func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting client streaming")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage %v", err)
	}

	req := []*calculatorpb.ComputeAverageRequest{
		{
			Number: 5,
		},
		{
			Number: 2,
		},
		{
			Number: 1,
		},
	}

	for _, v := range req {
		stream.Send(v)
		time.Sleep(500 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving from server %v", err)
	}
	fmt.Printf("ComputeAverage response: %v", res.Average)

}
func doBiDirectStreaming(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Starting BiDi streaming...")

	// we create a stream by invoking a client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream %v \n", err)
	}

	requests := []*calculatorpb.FindMaximumRequest{
		{
			Number: 1,
		},
		{
			Number: 3,
		},
		{
			Number: 2,
		},
		{
			Number: 1,
		},
		{
			Number: 5,
		},
		{
			Number: 20,
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
func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting SquareRoot unary (with error)...")

	// error call
	fmt.Println("Error request")
	fmt.Println()
	doErrorUnary_Call(c, int64(-5))

	fmt.Println(strings.Repeat("_", 25))

	// good call
	fmt.Println("Good request")
	fmt.Println()
	doErrorUnary_Call(c, int64(99))

}

func doErrorUnary_Call(c calculatorpb.CalculatorServiceClient, n int64) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})
	if err != nil {
		respErr, ok := status.FromError(err)

		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())

		} else {
			log.Fatalf("Big error calling SquareRoot %v", err)
		}
	}

	fmt.Printf("Answer: %v \n", res.GetNumberRoot())

}
