package main

import (
	"context"
	"fmt"
	"grpc-golang-master-class-build-modern-api-and-microservices/calculator/calculatorpb"
	"io"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I am a grpc client")

	certFile := "ssl/ca.crt"
	creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
	if sslErr != nil {
		log.Fatalf("Failed to load ssl files: %v", sslErr)
		return
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial: %v\n", err)
		return
	}
	defer conn.Close()

	client := calculatorpb.NewCalculatorServiceClient(conn)
	// fmt.Printf("Created grpc client: %f", client)

	// doUnary(client)

	// doServerStreaming(client)

	// doClientStreaming(client)

	// doBiDiStreaming(client)

	doSquareRootUnary(client)

	doSquareRootErrorUnary(client)
}

func doUnary(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a unary RPC...")

	req := &calculatorpb.SumRequest{
		FirstNumber:  1,
		SecondNumber: 3,
	}
	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call Sum function: %v\n", err)
	}
	fmt.Printf("Successfully called Sum function: %v\n", res)
}

func doServerStreaming(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a server streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 9999,
	}
	resStream, err := client.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call GreetManyTimes function: %v\n", err)
		return
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// reached end of stream
			break
		}
		if err != nil {
			log.Fatalf("Failed to read PrimeNumberDecomposition stream: %v", err)
			return
		}

		fmt.Printf("Successfully read PrimeNumberDecomposition stream: %v\n", msg.GetPrimeFactor())
	}

	fmt.Println("Successfully called PrimeNumberDecomposition function")
}

func randomNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 999
	return rand.Intn(max-min+1) + min
}

func doClientStreaming(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a client streaming RPC...")

	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Failed to call LongGreet function: %v\n", err)
		return
	}

	for i := 0; i < 10; i++ {
		req := &calculatorpb.ComputeAverageRequest{
			Number: int32(randomNumber()),
		}
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to close and receive ComputeAverage stream and response: %v\n", err)
		return
	}
	fmt.Printf("ComputeAverage response: %v\n", res.GetAverage())
}

func doBiDiStreaming(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a bi-directional streaming RPC...")

	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Failed to call FindMaximum function: %v\n", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		for i := 0; i < 10; i++ {
			req := &calculatorpb.FindMaximumRequest{
				Number: int32(randomNumber()),
			}
			fmt.Printf("Sending request: %v\n", req)
			stream.Send(req)
			time.Sleep(100 * time.Millisecond)
		}
		err = stream.CloseSend()
		if err != nil {
			log.Fatalf("Failed to close FindMaximum stream: %v\n", err)
			return
		}
	}()

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read FindMaximum stream: %v\n", err)
				break
			}

			fmt.Printf("Successfully read FindMaximum stream: %v\n", msg.GetMaximum())
		}
		close(waitc)
	}()

	<-waitc

	fmt.Println("FindMaximum successfully called")
}

func doSquareRootUnary(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a unary SquareRoot RPC...")

	req := &calculatorpb.SquareRootRequest{
		Number: 10,
	}
	res, err := client.SquareRoot(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call SquareRoot function: %v\n", err)
	}
	fmt.Printf("Successfully called SquareRoot function: %v\n", res)
}

func doSquareRootErrorUnary(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a unary SquareRootError RPC...")

	req := &calculatorpb.SquareRootRequest{
		Number: -1,
	}
	res, err := client.SquareRoot(context.Background(), req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Printf("Failed to call SquareRoot function Message: %v\n", respErr.Message())
			fmt.Printf("Failed to call SquareRoot function Code: %v\n", respErr.Code().String())
			log.Fatalf("Failed to call SquareRoot function Details: %v\n", respErr.Details()...)
			return
		}
		log.Fatalf("Failed to get SquareRoot grpc error: %v\n", err)
		return
	}

	fmt.Printf("Successfully called SquareRoot function: %v\n", res)
}
