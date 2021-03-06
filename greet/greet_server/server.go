package main

import (
	"context"
	"fmt"
	"grpc-golang-master-class-build-modern-api-and-microservices/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := strings.TrimSpace(req.GetGreeting().GetFirstName())
	lastName := strings.TrimSpace(req.GetGreeting().GetLastName())
	res := &greetpb.GreetResponse{
		Result: "Hello " + firstName + " " + lastName,
	}
	fmt.Println("Greet function was successful")
	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)
	firstName := strings.TrimSpace(req.GetGreeting().GetFirstName())
	lastName := strings.TrimSpace(req.GetGreeting().GetLastName())
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: "Hello " + firstName + " " + lastName + " number " + strconv.Itoa(i),
		}
		stream.Send(res)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("GreetManyTimes function was successful")
	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked")

	result := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// done reading stream
			er := stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
			if er != nil {
				log.Fatalf("Error while sending and closing client stream: %v\n", er)
				return err
			}
			fmt.Println("LongGreet function was successful")
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
			return err
		}

		firstName := strings.TrimSpace(req.GetGreeting().GetFirstName())
		lastName := strings.TrimSpace(req.GetGreeting().GetLastName())

		greeting := "Hello " + firstName + " " + lastName + "! "

		fmt.Printf("LongGreet function saying hi to %v\n", greeting)

		result += greeting
	}
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function was invoked")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("GreetEveryone function was successful")
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
			return err
		}

		firstName := strings.TrimSpace(req.GetGreeting().GetFirstName())
		lastName := strings.TrimSpace(req.GetGreeting().GetLastName())

		result := "Hello " + firstName + " " + lastName + "!"

		fmt.Printf("GreetEveryone function result %v\n", result)

		// done reading stream
		er := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if er != nil {
			log.Fatalf("Error while sending and closing client stream: %v\n", er)
			return err
		}
	}
}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked with %v\n", req)
	for i := 0; i < 4; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("client canceled request to GreetWithDeadline")
			return nil, status.Error(codes.Canceled, "client canceled request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := strings.TrimSpace(req.GetGreeting().GetFirstName())
	lastName := strings.TrimSpace(req.GetGreeting().GetLastName())
	res := &greetpb.GreetWithDeadlineResponse{
		Result: "Hello " + firstName + " " + lastName,
	}
	fmt.Println("GreetWithDeadline function was successful")
	return res, nil
}

func main() {
	fmt.Println("greet server")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env: %v", err)
		return
	}

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}

	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"

	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("Failed to load ssl files: %v", sslErr)
		return
	}

	s := grpc.NewServer(grpc.Creds(creds))
	greetpb.RegisterGreetServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		return
	}
}
