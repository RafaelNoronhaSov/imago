package main

import (
	"log"
	"net"

	"imago-tree/root"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	rootServer := root.NewServer()
	root.RegisterRootServiceServer(grpcServer, rootServer)

	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on %s", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %s: %v", port, err)
	}
}
