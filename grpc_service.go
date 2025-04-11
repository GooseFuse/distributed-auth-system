package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"sync"

	distributed_auth_system "github.com/GooseFuse/distributed-auth-system/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TransactionService defines the gRPC service for handling transactions
type TransactionService struct {
	distributed_auth_system.UnimplementedTransactionServiceServer
	dataStore        *DataStore
	consensusManager *ConsensusManager
	nodeID           string
	mutex            sync.RWMutex
}

// NewTransactionService initializes a new TransactionService
func NewTransactionService(nodeID string, dataStore *DataStore, consensusManager *ConsensusManager) *TransactionService {
	return &TransactionService{
		dataStore:        dataStore,
		consensusManager: consensusManager,
		nodeID:           nodeID,
	}
}

// StartServer starts the gRPC server
func StartServer(nodeID string, port string, dataStore *DataStore, consensusManager *ConsensusManager, useTLS bool) {
	var opts []grpc.ServerOption

	if useTLS {
		// Load TLS credentials
		cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
		if err != nil {
			log.Fatalf("Failed to load certificates: %v", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}

		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(opts...)
	service := NewTransactionService(nodeID, dataStore, consensusManager)
	distributed_auth_system.RegisterTransactionServiceServer(grpcServer, service)

	log.Printf("Node %s: gRPC server listening on %s", nodeID, port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// HandleTransaction handles a transaction request
func (s *TransactionService) HandleTransaction(ctx context.Context, req *distributed_auth_system.TransactionRequest) (*distributed_auth_system.TransactionResponse, error) {
	// Log the transaction
	log.Printf("Received transaction: key=%s, value=%s", req.Key, req.Value)

	// If this node is not the leader, forward to leader (in a real implementation)
	if !s.consensusManager.IsLeader() {
		log.Printf("Not the leader, would forward to leader in a real implementation")
		// In a real implementation, we would forward to the leader
		// For this example, we'll just handle it locally
	}

	// Propose transaction to the consensus mechanism
	success, err := s.consensusManager.ProposeTransaction(req.Key, req.Value)
	if err != nil {
		log.Printf("Error proposing transaction: %v", err)
		return &distributed_auth_system.TransactionResponse{
			Success: false,
		}, nil
	}

	return &distributed_auth_system.TransactionResponse{
		Success: success,
	}, nil
}
