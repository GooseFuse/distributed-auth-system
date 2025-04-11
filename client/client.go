// This file is a client example for the distributed-auth-system
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"time"

	distributed_auth_system "github.com/GooseFuse/distributed-auth-system/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Parse command line flags
	serverAddr := flag.String("server", "localhost:50051", "Server address")
	useTLS := flag.Bool("tls", false, "Use TLS")
	operation := flag.String("op", "store", "Operation: store, get, or auth")
	key := flag.String("key", "user123", "Key for store/get operations")
	value := flag.String("value", "password456", "Value for store operation")
	flag.Parse()

	// Set up connection options
	var opts []grpc.DialOption
	if *useTLS {
		// Load TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, // For testing only, don't use in production
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Connect to the server
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := distributed_auth_system.NewTransactionServiceClient(conn)

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
func storeData(ctx context.Context, client distributed_auth_system.TransactionServiceClient, key, value string) {
	fmt.Printf("Storing data: %s = %s\n", key, value)

	// Create a transaction request
	req := &distributed_auth_system.TransactionRequest{
		Key:   key,
		Value: value,
	}

	// Send the request
	resp, err := client.HandleTransaction(ctx, req)
	if err != nil {
		log.Fatalf("Failed to store data: %v", err)
	}

	// Print the response
	if resp.Success {
		fmt.Println("Data stored successfully")
	} else {
		// Note: In the current implementation, TransactionResponse doesn't have an ErrorMessage field
		fmt.Println("Failed to store data")
	}
}

// getData retrieves a value by key
func getData(ctx context.Context, client distributed_auth_system.TransactionServiceClient, key string) {
	fmt.Printf("Retrieving data for key: %s\n", key)

	// In a real implementation, we would have a separate RPC for getting data
	// For this example, we'll use the HandleTransaction RPC with an empty value
	// to indicate a get operation
	req := &distributed_auth_system.TransactionRequest{
		Key:   key,
		Value: "", // Empty value indicates a get operation
	}

	// Send the request
	resp, err := client.HandleTransaction(ctx, req)
	if err != nil {
		log.Fatalf("Failed to retrieve data: %v", err)
	}

	// Print the response
	if resp.Success {
		fmt.Println("Data retrieved successfully")
		// In a real implementation, the response would include the value
		fmt.Println("Note: This is a simplified example. In a real implementation, the response would include the value.")
	} else {
		// Note: In the current implementation, TransactionResponse doesn't have an ErrorMessage field
		fmt.Println("Failed to retrieve data")
	}
}

// authenticateUser simulates user authentication
func authenticateUser(ctx context.Context, client distributed_auth_system.TransactionServiceClient, username, password string) {
	fmt.Printf("Authenticating user: %s\n", username)

	// In a real implementation, we would have a separate AuthService
	// For this example, we'll use the HandleTransaction RPC with a special key format
	// to indicate an authentication operation
	authKey := fmt.Sprintf("auth:%s", username)

	// Store the password hash (in a real implementation, this would be done during user registration)
	storeReq := &distributed_auth_system.TransactionRequest{
		Key:   authKey,
		Value: password, // In a real implementation, this would be a password hash
	}

	storeResp, err := client.HandleTransaction(ctx, storeReq)
	if err != nil {
		log.Fatalf("Failed to store authentication data: %v", err)
	}

	if !storeResp.Success {
		// Note: In the current implementation, TransactionResponse doesn't have an ErrorMessage field
		fmt.Println("Failed to store authentication data")
		return
	}

	// Simulate authentication by checking if the stored password matches
	// In a real implementation, this would be a separate Login RPC
	authReq := &distributed_auth_system.TransactionRequest{
		Key:   authKey,
		Value: password,
	}

	authResp, err := client.HandleTransaction(ctx, authReq)
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	if authResp.Success {
		fmt.Println("Authentication successful")
		fmt.Println("Note: This is a simplified example. In a real implementation, the response would include an authentication token.")
	} else {
		// Note: In the current implementation, TransactionResponse doesn't have an ErrorMessage field
		fmt.Println("Authentication failed")
	}
}
