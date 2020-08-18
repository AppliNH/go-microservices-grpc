package calculator_client

import (
	"applinh/gogrpcudemy/calculator/calculatorpb"
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

	c := calculatorpb.NewCalculatorServiceClient(conn)

	doUnary(c)
	doServerStreaming(c)

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
