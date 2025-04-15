package node

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

type Node struct {
	dbPath            string
	nodesDir          string
	config            *Config
	raftConfig        *RaftConfig
	dataStore         *DataStore
	securityManager   *SecurityManager
	networkManager    *NetworkManager
	consensusManager  *ConsensusManager
	syncManager       *SyncManager
	checkpointManager *CheckpointManager
}

func newNode(dbPath string, config *Config, raftConfig *RaftConfig) *Node {
	return &Node{
		dbPath:     dbPath,
		config:     config,
		raftConfig: raftConfig,
	}
}

func Create() {
	// Parse OS vars
	nodeID := os.Getenv("NODE_ID")
	redisAddr := os.Getenv("REDIS_URL")
	port := os.Getenv("PORT")
	//mode := os.Getenv("MODE")

	// Parse command line flags
	dbPath := flag.String("db", "./data", "Path to database")
	useTLS := flag.Bool("tls", false, "Use TLS")
	flag.Parse()

	// Initialize configuration
	config := &Config{
		NodeID:             nodeID,
		Port:               port,
		RedisAddr:          redisAddr,
		UseTLS:             *useTLS,
		ElectionTimeoutMin: 15000, // 150ms
		ElectionTimeoutMax: 30000, // 300ms
		HeartbeatInterval:  50,    // 50ms
		PeerAddresses: map[string]string{
			"node1": "auth-node-1:6333",
			"node2": "auth-node-2:6333",
			"node3": "auth-node-3:6333",
		},
	}

	// Initialize the Raft configuration
	raftConfig := &RaftConfig{
		NodeID:             config.NodeID,
		PeerAddresses:      config.PeerAddresses,
		ElectionTimeoutMin: config.ElectionTimeoutMin,
		ElectionTimeoutMax: config.ElectionTimeoutMax,
		HeartbeatInterval:  config.HeartbeatInterval,
	}

	node := newNode(*dbPath, config, raftConfig)
	node.start()
	defer node.stop()

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Shutting down...")
	time.Sleep(1 * time.Second) // Give time for graceful shutdown
}

func (n *Node) start() {
	n.createNodesDir()
	n.createDataDir()

	var err error
	// Initialize the data store
	n.dataStore, err = NewDataStore(n.config.DBPath, n.config.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to initialize data store: %v", err)
	}

	// Initialize the security manager
	n.securityManager, err = NewSecurityManager(n.config.NodeID, n.getCertDir(), n.config.UseTLS)
	if err != nil {
		log.Fatalf("Failed to initialize security manager: %v", err)
	}

	// Initialize the network manager
	n.networkManager = NewNetworkManager(n.config.NodeID, n.securityManager)

	// Add peers to the network manager
	for peerID, peerAddr := range n.config.PeerAddresses {
		if peerID != n.config.NodeID {
			n.networkManager.AddPeer(peerID, peerAddr)
		}
	}
	// Start the network manager
	n.networkManager.Start()

	// Initialize the consensus manager
	n.consensusManager = NewConsensusManager(*n.raftConfig, n.dataStore, n.networkManager)
	n.consensusManager.Start()

	// Initialize the sync manager
	n.syncManager = NewSyncManager(n.dataStore, n.networkManager)
	n.syncManager.Start()

	// Initialize the checkpoint manager
	checkpointDir := filepath.Join(n.config.DBPath, "checkpoints")
	n.checkpointManager = NewCheckpointManager(n.dataStore, checkpointDir, 5*time.Minute)
	n.checkpointManager.Start()

	// Start the gRPC server in a separate goroutine
	go StartServer(n.config.NodeID, n.config.Port, n.dataStore, n.consensusManager, n.config.UseTLS, n.getCertDir())

	// Log startup information
	log.Printf("Node %s started successfully", n.config.NodeID)
	log.Printf("Data directory: %s", n.config.DBPath)
	log.Printf("Redis address: %s", n.config.RedisAddr)
	log.Printf("Listening on: %s", n.config.Port)
	log.Printf("TLS enabled: %v", n.config.UseTLS)
	log.Printf("Connected to %d peers", len(n.config.PeerAddresses)-1)
}

func (n *Node) getCertDir() string {
	return filepath.Join(n.config.DBPath, "certs")
}

func (n *Node) stop() {
	n.checkpointManager.Stop()
	n.syncManager.Stop()
	n.consensusManager.Stop()
	n.networkManager.Stop()
	n.dataStore.Close()
}

func (n *Node) createNodesDir() {
	// Create the nodes directory if it doesn't exist
	n.nodesDir = filepath.Join(n.dbPath, "nodes")
	err := os.MkdirAll(n.nodesDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create nodes directory: %v", err)
	}
}

func (n *Node) createDataDir() {
	// Create the node-specific directory
	n.config.DBPath = filepath.Join(n.nodesDir, n.config.NodeID)
	err := os.MkdirAll(n.config.DBPath, 0755)
	if err != nil {
		log.Fatalf("Failed to create node directory: %v", err)
	}
}
