package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// Configuration parameters
type Config struct {
	NodeID             string
	Port               string
	DBPath             string
	RedisAddr          string
	PeerAddresses      map[string]string
	UseTLS             bool
	ElectionTimeoutMin int
	ElectionTimeoutMax int
	HeartbeatInterval  int
}

func main() {
	// Parse command line flags
	nodeID := flag.String("node", "node1", "Node ID")
	port := flag.String("port", ":50051", "Port to listen on")
	dbPath := flag.String("db", "./data", "Path to database")
	redisAddr := flag.String("redis", "localhost:6379", "Redis address")
	useTLS := flag.Bool("tls", false, "Use TLS")
	flag.Parse()

	// Create the data directory if it doesn't exist
	dataDir := filepath.Join(*dbPath, *nodeID)
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize configuration
	config := &Config{
		NodeID:             *nodeID,
		Port:               *port,
		DBPath:             dataDir,
		RedisAddr:          *redisAddr,
		UseTLS:             *useTLS,
		ElectionTimeoutMin: 150, // 150ms
		ElectionTimeoutMax: 300, // 300ms
		HeartbeatInterval:  50,  // 50ms
		PeerAddresses: map[string]string{
			"node1": "localhost:50051",
			"node2": "localhost:50052",
			"node3": "localhost:50053",
		},
	}

	// Initialize the data store
	dataStore, err := NewDataStore(config.DBPath, config.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to initialize data store: %v", err)
	}
	defer dataStore.Close()

	// Initialize the security manager
	certDir := filepath.Join(config.DBPath, "certs")
	securityManager, err := NewSecurityManager(config.NodeID, certDir)
	if err != nil {
		log.Fatalf("Failed to initialize security manager: %v", err)
	}

	// Initialize the network manager
	networkManager := NewNetworkManager(config.NodeID, securityManager)

	// Add peers to the network manager
	for peerID, peerAddr := range config.PeerAddresses {
		if peerID != config.NodeID {
			networkManager.AddPeer(peerID, peerAddr)
		}
	}

	// Start the network manager
	networkManager.Start()
	defer networkManager.Stop()

	// Initialize the Raft configuration
	raftConfig := RaftConfig{
		NodeID:             config.NodeID,
		PeerAddresses:      config.PeerAddresses,
		ElectionTimeoutMin: config.ElectionTimeoutMin,
		ElectionTimeoutMax: config.ElectionTimeoutMax,
		HeartbeatInterval:  config.HeartbeatInterval,
	}

	// Initialize the consensus manager
	consensusManager := NewConsensusManager(raftConfig, dataStore)
	consensusManager.Start()
	defer consensusManager.Stop()

	// Initialize the sync manager
	syncManager := NewSyncManager(dataStore, networkManager)
	syncManager.Start()
	defer syncManager.Stop()

	// Initialize the checkpoint manager
	checkpointDir := filepath.Join(config.DBPath, "checkpoints")
	checkpointManager := NewCheckpointManager(dataStore, checkpointDir, 5*time.Minute)
	checkpointManager.Start()
	defer checkpointManager.Stop()

	// Start the gRPC server in a separate goroutine
	go StartServer(config.NodeID, config.Port, dataStore, consensusManager, config.UseTLS)

	// Log startup information
	log.Printf("Node %s started successfully", config.NodeID)
	log.Printf("Data directory: %s", config.DBPath)
	log.Printf("Redis address: %s", config.RedisAddr)
	log.Printf("Listening on: %s", config.Port)
	log.Printf("TLS enabled: %v", config.UseTLS)
	log.Printf("Connected to %d peers", len(config.PeerAddresses)-1)

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Shutting down...")
	time.Sleep(1 * time.Second) // Give time for graceful shutdown
}
