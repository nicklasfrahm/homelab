package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/nicklasfrahm/cloud/pkg/blobkv"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"google.golang.org/grpc"
)

var (
	listen = flag.String("listen", ":2379", "Address to listen on")
)

func main() {
	flag.Parse()

	// Load .env file if it exists
	godotenv.Load()

	blobURL := os.Getenv("BLOB_URI")
	if blobURL == "" {
		log.Fatal("BLOB_URI environment variable is required")
	}

	store, err := blobkv.NewStore(blobURL)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	server := grpc.NewServer()
	etcdserverpb.RegisterKVServer(server, blobkv.NewKVServer(store))
	etcdserverpb.RegisterWatchServer(server, blobkv.NewWatchServer(store))

	listener, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting server on %s", *listen)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
