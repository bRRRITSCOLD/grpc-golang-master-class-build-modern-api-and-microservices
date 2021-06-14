package main

import (
	"context"
	"fmt"
	"grpc-golang-master-class-build-modern-api-and-microservices/greet/greetpb"
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

	client := greetpb.NewGreetServiceClient(conn)
	// fmt.Printf("Created grpc client: %f", client)

	doUnary(client)

	doServerStreaming(client)
}

func doUnary(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a unary RPC...")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bobby",
			LastName:  "Bushay",
		},
	}
	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call Greet function: %v", err)
	}
	fmt.Printf("Successfully called Greet function: %v\n", res)
}

func doServerStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a server streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bobby",
			LastName:  "Bushay",
		},
	}
	resStream, err := client.GreetManyTimes(context.Background(), req)
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
			log.Fatalf("Failed to read GreetManyTimes stream: %v", err)
		}

		fmt.Printf("Successfully read GreetManyTimes stream: %v\n", msg.GetResult())
	}

	fmt.Println("Successfully called GreetManyTimes function")
}
