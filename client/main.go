package main

import (
	"context"
	"fmt"
	"log"
	"time"

	timepb "github.com/gherlein/time-services/time-service-pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := timepb.NewTimeServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetTime(ctx, &timepb.TimeRequest{})
	if err != nil {
		log.Fatalf("Error calling GetTime: %v", err)
	}
	if res != nil {
		fmt.Printf("Current Time: %s\n", res.GetCurrentTime())
	}

}
