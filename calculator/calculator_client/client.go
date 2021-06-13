package main

import (
	"context"
	"fmt"
	"grpc-golang-master-class-build-modern-api-and-microservices/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I am a grpc client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := calculatorpb.NewCalculatorServiceClient(conn)
	// fmt.Printf("Created grpc client: %f", client)

	doUnary(client)

	doServerStreaming(client)
}

func doUnary(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a unary RPC...")

	req := &calculatorpb.SumRequest{
		FirstNumber:  1,
		SecondNumber: 3,
	}
	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call Sum function: %v", err)
	}
	fmt.Printf("Successfully called Sum function: %v", res)
}

func doServerStreaming(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a server streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 9999,
	}
	resStream, err := client.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call GreetManyTimes function: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// reached end of stream
			break
		}
		if err != nil {
			log.Fatalf("Failed to read PrimeNumberDecomposition stream: %v", err)
		}

		fmt.Printf("Successfully read PrimeNumberDecomposition stream: %v\n", msg.GetPrimeFactor())
	}

	fmt.Println("Successfully called PrimeNumberDecomposition function")
}
