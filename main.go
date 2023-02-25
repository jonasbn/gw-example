package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/indiependente/gw-example/rpc"
	gw "github.com/indiependente/gw-example/rpc/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
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
	reflection.Register(grpcServer)

	srv := &rpc.MsgAPISrv{
		Log: log.New(os.Stdout, ">>> ", 0),
	}

	gw.RegisterMessageAPIServiceServer(grpcServer, srv)

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
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterMessageAPIServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	srv := http.Server{
		Addr:        ":8080",
		ReadTimeout: 10 * time.Second,
		Handler:     mux,
	}

	return srv.ListenAndServe()
}

func main() {
	flag.Parse()

	go func() {
		log.Fatal(startGRPC())
	}()

	log.Fatal(startGW())
}
