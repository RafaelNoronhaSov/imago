// main.go
package main

import (
	"log"
	"net"

	"imago-node/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// main starts the gRPC server for the node service on port 50052.
//
// It listens on the specified port, creates a new gRPC server, creates a new
// node server, registers the node service with the server, enables reflection
// for debugging, logs the listening port, and serves the server over the
// specified port.
func main() {
	port := ":50052"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	nodeServer := server.NewServer()

	server.RegisterNodeServiceServer(grpcServer, nodeServer)

	reflection.Register(grpcServer)

	log.Printf("Node gRPC server listening on %s", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %s: %v", port, err)
	}
}
