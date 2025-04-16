package raft

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"time"

	"github.com/GooseFuse/distributed-auth-system/protoc"
	"github.com/GooseFuse/distributed-auth-system/system/datastore"
	"github.com/willf/bloom"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// TransactionService defines the gRPC service for handling transactions
type TransactionService struct {
	protoc.UnimplementedTransactionServiceServer
	dataStore        *datastore.DataStore
	consensusManager *ConsensusManager
	nodeID           string
	//mutex            sync.RWMutex
}

// NewTransactionService initializes a new TransactionService
func NewTransactionService(nodeID string, dataStore *datastore.DataStore, consensusManager *ConsensusManager) protoc.TransactionServiceServer {
	return &TransactionService{
		dataStore:        dataStore,
		consensusManager: consensusManager,
		nodeID:           nodeID,
	}
}

// StartServer starts the gRPC server
func StartServer(nodeID string, port string, dataStore *datastore.DataStore, consensusManager *ConsensusManager, useTLS bool, certDir string) {
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

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", nodeID, port))
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
		fmt.Printf("Not the leader, forward to leader...")
		leader := s.consensusManager.GetLeader()
		if leader == "" {
			fmt.Printf("Leader is not selected yet, dropping...")
			return &protoc.TransactionResponse{
				Success: false,
			}, nil
		} else {
			fmt.Printf("Not the leader. Please retry with leader: %s", leader)
			return &protoc.TransactionResponse{
				Success:      false,
				ErrorMessage: fmt.Sprintf("Not the leader. Please retry with leader: %s", leader),
				LeaderUrl:    s.consensusManager.raftNode.networkManager.GetPeerUrl(leader),
			}, nil
		}
	}

	// Propose transaction to the consensus mechanism
	log.Printf("Propose transaction: key=%s, value=%s", req.Transaction.Key, req.Transaction.Value)
	success, err := s.consensusManager.ProposeTransaction(req.Transaction.Key, req.Transaction.Value)
	if err != nil {
		log.Printf("Error proposing transaction: %v", err)
		return &protoc.TransactionResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &protoc.TransactionResponse{
		Success: success,
	}, nil
}

func (s *TransactionService) VerifyState(ctx context.Context, req *protoc.StateVerificationRequest) (svr *protoc.StateVerificationResponse, e error) {
	localMerkleRoot := s.dataStore.GetMerkleRoot() // ← your local copy
	log.Printf("🔁 VerifyState request from node %s | Peer Merkle: %s | Local Merkle: %s",
		req.NodeId, req.MerkleRoot, localMerkleRoot)

	isConsistent := req.MerkleRoot == localMerkleRoot

	return &protoc.StateVerificationResponse{
		Consistent: isConsistent,
		MerkleRoot: localMerkleRoot,
	}, nil
}

func (s *TransactionService) AppendEntries(ctx context.Context, req *protoc.AppendEntriesRequest) (*protoc.AppendEntriesResponse, error) {
	node := s.consensusManager.raftNode

	node.mutex.Lock()
	defer node.mutex.Unlock()

	// Step 1: Reject if term is outdated
	if int(req.Term) < node.currentTerm {
		return &protoc.AppendEntriesResponse{
			Term:    int32(node.currentTerm),
			Success: false,
		}, nil
	}

	// Step 2: If term is higher, update local term and step down if leader
	if int(req.Term) > node.currentTerm {
		node.currentTerm = int(req.Term)
		node.state = Follower
		node.votedFor = ""
	}

	// Step 3: Reset election timeout
	node.leaderID = req.LeaderId
	node.lastHeartbeat = time.Now()

	// Step 4: Validate log consistency
	if int(req.PrevLogIndex) >= 0 {
		if int(req.PrevLogIndex) >= len(node.log) ||
			node.log[req.PrevLogIndex].Term != int32(req.PrevLogTerm) {
			return &protoc.AppendEntriesResponse{
				Term:    int32(node.currentTerm),
				Success: false,
			}, nil
		}
	}

	// Step 5: Delete conflicting entries and append new ones
	for i, entry := range req.Entries {
		logIndex := int(req.PrevLogIndex) + 1 + i
		if logIndex < len(node.log) {
			if node.log[logIndex].Term != int32(entry.Term) {
				node.log = node.log[:logIndex] // delete conflicting entries
				break
			}
		}
	}

	for i, entry := range req.Entries {
		logIndex := int(req.PrevLogIndex) + 1 + i
		if logIndex >= len(node.log) {
			node.log = append(node.log, entry)
		}
	}

	// Step 6: Update commit index and apply new entries
	if int(req.LeaderCommit) > node.commitIndex {
		node.commitIndex = min(int(req.LeaderCommit), len(node.log)-1)

		for i := node.lastApplied + 1; i <= node.commitIndex; i++ {
			node.applyCh <- node.log[i]
			node.lastApplied = i
		}
	}

	return &protoc.AppendEntriesResponse{
		Term:    int32(node.currentTerm),
		Success: true,
	}, nil
}

func (s *TransactionService) SyncState(ctx context.Context, req *protoc.SyncRequest) (sr *protoc.SyncResponse, e error) {
	localMerkleRoot := s.dataStore.GetMerkleRoot()

	log.Printf("🔁 SyncState request from node %s | Peer Merkle: %s | Local Merkle: %s",
		req.NodeId, req.MerkleRoot, localMerkleRoot)

	// If Merkle roots are equal, no sync needed
	if req.MerkleRoot == localMerkleRoot {
		return &protoc.SyncResponse{
			Success:    true,
			MerkleRoot: localMerkleRoot,
			// No missing transactions
		}, nil
	}

	var missing []*protoc.Transaction

	allTxs, err := s.dataStore.GetAllTransactions()
	if err != nil {
		return nil, err
	}

	for _, tx := range allTxs {
		// Use bloom filter to avoid sending keys peer already has
		if !mightContainKey(req.BloomFilter, tx.Key) {
			missing = append(missing, tx)
		}
	}

	log.Printf("📤 Sending %d missing transactions to node %s", len(missing), req.NodeId)

	return &protoc.SyncResponse{
		Success:             true,
		MerkleRoot:          localMerkleRoot,
		MissingTransactions: missing,
	}, nil
}

func (s *TransactionService) RequestVote(ctx context.Context, in *protoc.RequestVoteRequest) (*protoc.RequestVoteResponse, error) {
	candidateID := in.CandidateId
	term := int(in.Term)
	lastLogIndex := int(in.LastLogIndex)
	lastLogTerm := int(in.LastLogTerm)
	currentTerm, voteGranted := s.consensusManager.raftNode.HanldeRequestVote(candidateID, term, lastLogIndex, lastLogTerm)
	return &protoc.RequestVoteResponse{Term: int32(currentTerm), VoteGranted: voteGranted}, nil
}

func mightContainKey(filterData []byte, key string) bool {
	filter := bloom.New(1, 1) // dummy values, will be overwritten
	err := filter.GobDecode(filterData)
	if err != nil {
		// fallback: if we can't decode, assume key might be present
		return true
	}
	return filter.Test([]byte(key))
}

func (s *TransactionService) Get(ctx context.Context, req *protoc.GetRequest) (*protoc.GetResponse, error) {
	// Log the transaction
	log.Printf("Received GetRequest: User=%s", req.User)

	// Validate the request
	if req.GetUser() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user is required")
	}

	// Fetch from your underlying store
	value, err := s.dataStore.GetData(req.GetUser())
	if err != nil {
		log.Printf("store get error: %v", err)
		return nil, status.Errorf(codes.Internal, "store get error: %v", err)
	}

	log.Printf("GetResponse: Value=%s", value)
	// Return the value (even if not found, value can be empty)
	return &protoc.GetResponse{
		Value:   value,
		Success: true,
	}, nil
}

func (s *TransactionService) Auth(ctx context.Context, req *protoc.AuthRequest) (*protoc.AuthResponse, error) {
	log.Printf("Received AuthRequest: User=%s Pass:%s", req.User, req.Pass)

	// Validate the request
	if req.GetUser() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "User is required")
	}
	if req.GetPass() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Pass is required")
	}

	// Fetch from your underlying store
	value, err := s.dataStore.GetData(req.GetUser())
	if err != nil {
		log.Printf("store get error: %v", err)
		return nil, status.Errorf(codes.Internal, "store get error: %v", err)
	}

	if value != req.GetPass() {
		log.Printf("Auth failed")
		return &protoc.AuthResponse{Success: false, ErrorMessage: "Invalid credentials"}, nil
	}

	log.Printf("Auth success")
	return &protoc.AuthResponse{
		Success: true,
	}, nil
}
