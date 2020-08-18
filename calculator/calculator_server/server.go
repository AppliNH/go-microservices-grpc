package calculator_server

import (
	"applinh/gogrpcudemy/calculator/calculatorpb"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

// from greet/greetpb/greet.pb.go - GreetServiceServer interface
func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Printf("Sum invoked with %v \n", req)
	a := req.GetA()
	b := req.GetB()

	result := a + b

	res := &calculatorpb.SumResponse{
		Result: result,
	}

	return res, nil

}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Printf("Greet stream invoked with %v \n", req)

	number := req.GetNumber()

	var k int64 = 2

	for number > 1 {

		if number%k == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeNumber: k,
			}
			stream.Send(res)
			number = number / k
		} else {
			k = k + 1
		}

	}
	return nil
}

func StartServer() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen %v \n", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
