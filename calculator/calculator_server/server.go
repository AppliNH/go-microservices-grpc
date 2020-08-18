package calculator_server

import (
	"applinh/gogrpcudemy/calculator/calculatorpb"
	"context"
	"io"
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
