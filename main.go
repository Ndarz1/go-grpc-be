package main

import (
	`log`
	`net`
	
	`google.golang.org/grpc`
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	
	serv := grpc.NewServer()
	
	log.Println("Server is running on :50051 port")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}
