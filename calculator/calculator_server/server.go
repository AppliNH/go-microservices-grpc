package calculator_server

import (
	"applinh/gogrpcudemy/calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	log.Printf("ComputeAverage stream invoked with a streaming request \n")
	numbers := []int64{}
	var average float64
	var sum int64
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// finished reading client stream

			for _, n := range numbers {
				sum += n
			}
			average = float64(sum) / float64(int64(len(numbers)))
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v \n", err)
		}
		numbers = append(numbers, msg.GetNumber())

	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	log.Printf("FindMaximum stream invoked with a bi-di streaming request \n")

	var maxNumber int64

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error occured on reading clientStream %v", err)
			return err
		}
		number := req.GetNumber()
		if number > maxNumber {
			maxNumber = number
		}

		err = stream.Send(&calculatorpb.FindMaximumResponse{
			Result: maxNumber,
		})
		if err != nil {
			log.Fatalf("Error when streaming data to client %v", err)
			return err
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("SquareRoot unary")

	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received a negative number: %v", number))
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

//-------------------------------------------
func StartServer() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen %v \n", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	reflection.Register(s) // allows us to expose the gRPC server so the client can see what's available. You can use Evans CLI for that

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
