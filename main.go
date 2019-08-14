package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	gw "github.com/indiependente/gw-example/rpc"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", ":9090", "gRPC server endpoint")
)

func startGRPC() error {
	log.Printf("Starting GRPC server on port %s...\n", *grpcServerEndpoint)
	lis, err := net.Listen("tcp", *grpcServerEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	srv := &gw.MsgAPISrv{
		Log: log.New(os.Stdout, ">>> ", 0),
	}

	gw.RegisterMessageAPIServer(grpcServer, srv)

	return grpcServer.Serve(lis)
}

func startGW() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log.Println("Starting HTTP server on port :8080...")
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterMessageAPIHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8080", mux)
}

func main() {
	flag.Parse()

	go func() {
		log.Fatal(startGRPC())
	}()

	log.Fatal(startGW())
}
