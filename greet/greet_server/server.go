package greet_server

import (
	"applinh/gogrpcudemy/greet/greetpb"
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

// from greet/greetpb/greet.pb.go - GreetServiceServer interface
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("Greet invoked with %v \n", req)
	first_name := req.GetGreeting().GetFirstName()
	last_name := req.GetGreeting().GetLastName()

	result := "Hello " + first_name + " " + last_name

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil

}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Printf("Greet stream invoked with %v \n", req)

	first_name := req.GetGreeting().GetFirstName()
	last_name := req.GetGreeting().GetLastName()

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: strconv.Itoa(i) + " Hello " + first_name + " " + last_name,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond) // Sleep 1 second
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Printf("LongGreet stream invoked with a streaming request \n")
	result := "Hello "
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// finished reading client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v \n", err)
		}
		first_name := msg.GetGreeting().GetFirstName()
		last_name := msg.GetGreeting().GetLastName()

		result += first_name + " " + last_name + ", "

	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Printf("GreetEveryone stream invoked with a bi-di streaming request \n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error occured on reading clientStream %v", err)
			return err
		}
		first_name := req.GetGreeting().GetFirstName()
		last_name := req.GetGreeting().GetLastName()

		result := "Hello " + first_name + " " + last_name + " !"

		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error when streaming data to client %v", err)
			return err
		}
	}

}

func StartServer() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen %v \n", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
