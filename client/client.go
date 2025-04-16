package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/GooseFuse/distributed-auth-system/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	useTLS *bool
)

func getConn(serverAddr string, useTLS bool) *grpc.ClientConn {
	// Set up connection options
	var opts []grpc.DialOption
	if useTLS {
		// Load TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, // For testing only, don't use in production
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Connect to the server
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	return conn
}

func main() {
	// Parse command line flags
	serverAddr := flag.String("server", "localhost:50051", "Server address")
	useTLS = flag.Bool("tls", false, "Use TLS")
	operation := flag.String("op", "store", "Operation: store, get, or auth")
	key := flag.String("key", "user123", "Key for store/get operations")
	value := flag.String("value", "password456", "Value for store operation")
	flag.Parse()

	conn := getConn(*serverAddr, *useTLS)
	defer conn.Close()

	// Create a client
	client := protoc.NewTransactionServiceClient(conn)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Perform the requested operation
	switch *operation {
	case "store":
		storeData(ctx, client, *key, *value)
	case "get":
		getData(ctx, client, *key)
	case "auth":
		authenticateUser(ctx, client, *key, *value)
	default:
		log.Fatalf("Unknown operation: %s", *operation)
	}
}

// storeData stores a key-value pair
func storeData(ctx context.Context, client protoc.TransactionServiceClient, key, value string) {
	fmt.Printf("Storing data: %s = %s\n", key, value)

	// Create a transaction request
	req := &protoc.TransactionRequest{
		Transaction: &protoc.Transaction{
			Key:   key,
			Value: value,
		},
	}

	// Send the request
	resp, err := client.HandleTransaction(ctx, req)
	if err != nil {
		log.Fatalf("Failed to store data: %v", err)
	}

	// Print the response
	if resp.Success {
		fmt.Println("Data stored successfully")
	} else if resp.LeaderUrl != "" {
		fmt.Printf("Redirected to: %s\n", resp.LeaderUrl)
		conn := getConn(resp.LeaderUrl, *useTLS)
		defer conn.Close()
		// Create a client
		leaderClient := protoc.NewTransactionServiceClient(conn)
		storeData(ctx, leaderClient, key, value)
	} else {
		fmt.Println("Failed to store data")
	}
}

// getData retrieves a value by key
func getData(ctx context.Context, client protoc.TransactionServiceClient, key string) {
	fmt.Printf("Retrieving data for key: %s\n", key)

	// In a real implementation, we would have a separate RPC for getting data
	// For this example, we'll use the HandleTransaction RPC with an empty value
	// to indicate a get operation
	req := &protoc.GetRequest{
		User: key,
	}

	// Send the request
	resp, err := client.Get(ctx, req)
	if err != nil {
		log.Fatalf("Failed to retrieve data: %v", err)
	}

	// Print the response
	if resp.Success {
		log.Printf("Data retrieved successfully: %s\n", resp.Value)
	} else {
		log.Printf("Failed to retrieve data: %s\n", resp.ErrorMessage)
	}
}

// authenticateUser simulates user authentication
func authenticateUser(ctx context.Context, client protoc.TransactionServiceClient, username, password string) {
	fmt.Printf("Authenticating user: %s\n", username)

	authReq := &protoc.AuthRequest{
		User: username,
		Pass: password,
	}

	authResp, err := client.Auth(ctx, authReq)
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	if !authResp.Success {
		log.Println("Failed to authenticate")
		return
	}

	log.Println("Authentication successful")
}
