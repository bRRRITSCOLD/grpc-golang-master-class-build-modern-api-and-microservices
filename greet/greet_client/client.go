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
		log.Fatalf("Failed to dial: %v", err)
		return
	}
	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)
	// fmt.Printf("Created grpc client: %f", client)

	doUnary(client)

	// doServerStreaming(client)

	// doClientStreaming(client)

	// doBiDiStreaming(client)

	// doUnaryGreetWithDeadline(client, 5*time.Second)
	// doUnaryGreetWithDeadline(client, 1*time.Second)
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

func doUnaryGreetWithDeadline(client greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a unary RPC...")

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bobby",
			LastName:  "Bushay",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.GreetWithDeadline(ctx, req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Printf("Failed to call GreetWithDeadline function Message: %v\n", respErr.Message())
			fmt.Printf("Failed to call GreetWithDeadline function Code: %v\n", respErr.Code().String())
			log.Fatalf("Failed to call GreetWithDeadline function Details: %v\n", respErr.Details()...)
			return
		}
		log.Fatalf("Failed to get GreetWithDeadline grpc error: %v\n", err)
		return
	}

	fmt.Printf("Successfully called GreetWithDeadline function: %v\n", res)
}
