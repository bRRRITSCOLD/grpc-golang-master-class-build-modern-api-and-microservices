package main

import (
	"context"
	"fmt"
	"grpc-golang-master-class-build-modern-api-and-microservices/greet/greetpb"
	"io"
	"log"
	"strconv"
	"time"

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

	doClientStreaming(client)

	doBiDiStreaming(client)
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
		log.Fatalf("Failed to call Greet function: %v\n", err)
		return
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
			log.Fatalf("Failed to read GreetManyTimes stream: %v\n", err)
			return
		}

		fmt.Printf("Successfully read GreetManyTimes stream: %v\n", msg.GetResult())
	}

	fmt.Println("Successfully called GreetManyTimes function")
}

func doClientStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a client streaming RPC...")

	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Failed to call LongGreet function: %v\n", err)
		return
	}

	for i := 0; i < 10; i++ {
		req := &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bobby " + strconv.Itoa(i+1),
				LastName:  "Bushay " + strconv.Itoa(i+1),
			},
		}
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to close and receive LongGreet stream and response: %v\n", err)
		return
	}
	fmt.Printf("LongGreet response: %v\n", res.GetResult())
}

func doBiDiStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a bi-directional streaming RPC...")

	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Failed to call GreetEveryone function: %v\n", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		for i := 0; i < 10; i++ {
			req := &greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Bobby " + strconv.Itoa(i+1),
					LastName:  "Bushay " + strconv.Itoa(i+1),
				},
			}
			fmt.Printf("Sending request: %v\n", req)
			stream.Send(req)
			time.Sleep(100 * time.Millisecond)
		}
		err = stream.CloseSend()
		if err != nil {
			log.Fatalf("Failed to close GreetEveryone stream: %v\n", err)
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
				log.Fatalf("Failed to read GreetManyTimes stream: %v\n", err)
				break
			}

			fmt.Printf("Successfully read GreetManyTimes stream: %v\n", msg.GetResult())
		}
		close(waitc)
	}()

	<-waitc

	fmt.Println("GreetEveryone successfully called")
}
