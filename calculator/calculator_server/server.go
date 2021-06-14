package main

import (
	"context"
	"fmt"
	"grpc-golang-master-class-build-modern-api-and-microservices/calculator/calculatorpb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()
	res := &calculatorpb.SumResponse{
		Result: firstNumber + secondNumber,
	}
	return res, nil
}

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)

	number := req.GetNumber()
	divisor := int32(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("PrimeNumberDecomposition divisor has increased %v\n", divisor)
		}
		time.Sleep(2 * time.Millisecond)
	}
	return nil
}

func (s *server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage function was invoked")

	total := 0
	amountOfNumbers := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// done reading stream
			fmt.Printf("ComputeAverage function total %v\n", total)
			fmt.Printf("ComputeAverage function amountOfNumbers %v\n", amountOfNumbers)
			er := stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: int32(total / amountOfNumbers),
			})
			if er != nil {
				log.Fatalf("Error while sending and closing client stream: %v\n", er)
				return err
			}
			fmt.Println("ComputeAverage function was successful")
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
			return err
		}

		number := req.GetNumber()
		total += int(number)
		amountOfNumbers += 1

		fmt.Printf("ComputeAverage function adding number %v\n", number)
	}
}

func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum function was invoked")

	maximum := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("FindMaximum function was successful")
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
			return err
		}

		number := req.GetNumber()

		fmt.Printf("FindMaximum received number %v\n", number)

		if number > maximum {
			maximum = number
		}

		fmt.Printf("FindMaximum maximum %v\n", maximum)

		// done reading stream
		er := stream.Send(&calculatorpb.FindMaximumResponse{
			Maximum: maximum,
		})
		if er != nil {
			log.Fatalf("Error while sending and closing client stream: %v\n", er)
			return err
		}
	}
}

func main() {
	fmt.Println("calculator server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
