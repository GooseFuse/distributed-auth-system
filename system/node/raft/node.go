package raft

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/GooseFuse/distributed-auth-system/system/datastore"
	"github.com/GooseFuse/distributed-auth-system/system/network"
	"github.com/GooseFuse/distributed-auth-system/system/security"
)

type Node struct {
	UseTLS            bool
	DBPath            string
	Port              string
	raftConfig        *RaftConfig
	consensusManager  *ConsensusManager
	dataStore         *datastore.DataStore
	networkManager    *network.NetworkManager
	securityManager   *security.SecurityManager
	syncManager       *SyncManager
	checkpointManager *CheckpointManager
}

func NewNode(NodeID string, PeerAddresses map[string]string, ElectionTimeoutMin int, ElectionTimeoutMax int, HeartbeatInterval int, dataStore *datastore.DataStore, UseTLS bool, dbPath string, port string) *Node {
	// Initialize the Raft configuration
	return &Node{
		DBPath:    dbPath,
		dataStore: dataStore,
		UseTLS:    UseTLS,
		Port:      port,
		raftConfig: &RaftConfig{
			NodeID:             NodeID,
			PeerAddresses:      PeerAddresses,
			ElectionTimeoutMin: ElectionTimeoutMin,
			ElectionTimeoutMax: ElectionTimeoutMax,
			HeartbeatInterval:  HeartbeatInterval,
		},
	}
}

func (n *Node) Start() {
	var err error
	redisAddr := os.Getenv("REDIS_URL")
	// Initialize the data store
	n.dataStore, err = datastore.NewDataStore(n.DBPath, redisAddr)
	if err != nil {
		log.Fatalf("Failed to initialize data store: %v", err)
	}

	// Initialize the security manager
	n.securityManager, err = security.NewSecurityManager(n.raftConfig.NodeID, n.getCertDir(), n.UseTLS)
	if err != nil {
		log.Fatalf("Failed to initialize security manager: %v", err)
	}

	// Initialize the network manager
	n.networkManager = network.NewNetworkManager(n.raftConfig.NodeID, n.securityManager)

	// Add peers to the network manager
	for peerID, peerAddr := range n.raftConfig.PeerAddresses {
		if peerID != n.raftConfig.NodeID {
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
	checkpointDir := filepath.Join(n.DBPath, "checkpoints")
	n.checkpointManager = NewCheckpointManager(n.dataStore, checkpointDir, 5*time.Minute)
	n.checkpointManager.Start()

	// Start the gRPC server in a separate goroutine
	go StartServer(n.raftConfig.NodeID, n.Port, n.dataStore, n.consensusManager, n.UseTLS, n.getCertDir())
}

func (n *Node) Stop() {
	n.consensusManager.Stop()
}

func (n *Node) getCertDir() string {
	return filepath.Join(n.DBPath, "certs")
}
