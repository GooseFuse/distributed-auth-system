package node

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"path/filepath"
	"sync"

	"github.com/GooseFuse/distributed-auth-system/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// TransactionService defines the gRPC service for handling transactions
type TransactionService struct {
	protoc.UnimplementedTransactionServiceServer
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
func StartServer(nodeID string, port string, dataStore *DataStore, consensusManager *ConsensusManager, useTLS bool, certDir string) {
	var opts []grpc.ServerOption

	if useTLS {
		// Load TLS credentials
		certPath := filepath.Join(certDir, "cert.pem")
		keyPath := filepath.Join(certDir, "key.pem")
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
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
	protoc.RegisterTransactionServiceServer(grpcServer, service)
	for svc := range grpcServer.GetServiceInfo() {
		log.Printf("✔ Registered gRPC service: %s", svc)
	}
	reflection.Register(grpcServer)
	log.Printf("Node %s: gRPC server listening on %s", nodeID, port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// HandleTransaction handles a transaction request
func (s *TransactionService) HandleTransaction(ctx context.Context, req *protoc.TransactionRequest) (*protoc.TransactionResponse, error) {
	// Log the transaction
	log.Printf("Received transaction: key=%s, value=%s", req.Transaction.Key, req.Transaction.Value)

	// If this node is not the leader, forward to leader (in a real implementation)
	if !s.consensusManager.IsLeader() {
		log.Printf("Not the leader, would forward to leader in a real implementation")
		// In a real implementation, we would forward to the leader
		// For this example, we'll just handle it locally
	}

	// Propose transaction to the consensus mechanism
	success, err := s.consensusManager.ProposeTransaction(req.Transaction.Key, req.Transaction.Value)
	if err != nil {
		log.Printf("Error proposing transaction: %v", err)
		return &protoc.TransactionResponse{
			Success: false,
		}, nil
	}

	return &protoc.TransactionResponse{
		Success: success,
	}, nil
}

func (s *TransactionService) VerifyState(ctx context.Context, req *protoc.StateVerificationRequest) (svr *protoc.StateVerificationResponse, e error) {
	localMerkleRoot := s.dataStore.GetMerkleRoot() // ← your local copy

	isConsistent := req.MerkleRoot == localMerkleRoot

	return &protoc.StateVerificationResponse{
		Consistent: isConsistent,
		MerkleRoot: localMerkleRoot,
	}, nil
}
