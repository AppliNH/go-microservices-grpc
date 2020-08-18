package calculator_client

import (
	"applinh/gogrpcudemy/calculator/calculatorpb"
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

	c := calculatorpb.NewCalculatorServiceClient(conn)

	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)

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
