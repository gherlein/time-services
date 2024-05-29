package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	timepb "github.com/gherlein/time-services/time-service-pb"
	"google.golang.org/grpc"
)

type server struct {
	timepb.UnimplementedTimeServiceServer
}

func (s *server) GetTime(ctx context.Context, req *timepb.TimeRequest) (*timepb.TimeResponse, error) {
	currentTime := time.Now().Format(time.RFC3339)
	println("Current time: ", currentTime)
	res := &timepb.TimeResponse{
		CurrentTime: currentTime, // ERROR: ChatGPT used Current_time
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	timepb.RegisterTimeServiceServer(s, &server{})

	fmt.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
