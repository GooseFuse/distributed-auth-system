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

	"github.com/GooseFuse/distributed-auth-system/system/datastore"
	"github.com/GooseFuse/distributed-auth-system/system/node/blockchain"
	"github.com/GooseFuse/distributed-auth-system/system/node/raft"
)

type Node struct {
	imp       NodeInterface
	dbPath    string
	nodesDir  string
	config    *Config
	dataStore *datastore.DataStore
}

type NodeType string

const (
	raftType       = "RAFT"
	blockchainType = "BLOCKCHAIN"
)

func newNode(mode NodeType, dbPath string, config *Config) *Node {
	n := &Node{
		dbPath: dbPath,
		config: config,
	}

	switch mode {
	case raftType:
		n.imp = raft.NewNode(config.NodeID, config.PeerAddresses, config.ElectionTimeoutMin, config.ElectionTimeoutMax, config.HeartbeatInterval, n.dataStore, config.UseTLS, dbPath, config.Port)
	case blockchainType:
		n.imp = &blockchain.Node{}

	}
	return n
}

func Create() {
	// Parse OS vars
	nodeID := os.Getenv("NODE_ID")
	redisAddr := os.Getenv("REDIS_URL")
	port := os.Getenv("PORT")
	mode := NodeType(os.Getenv("MODE"))

	if mode == "" {
		mode = raftType
	}

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

	node := newNode(mode, *dbPath, config)
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

	n.imp.Start()

	// Log startup information
	log.Printf("Node %s started successfully", n.config.NodeID)
	log.Printf("Data directory: %s", n.config.DBPath)
	log.Printf("Redis address: %s", n.config.RedisAddr)
	log.Printf("Listening on: %s", n.config.Port)
	log.Printf("TLS enabled: %v", n.config.UseTLS)
	log.Printf("Connected to %d peers", len(n.config.PeerAddresses)-1)
}

func (n *Node) stop() {
	n.imp.Stop()
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
