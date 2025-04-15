package raft

// RaftConfig contains configuration parameters for the Raft consensus algorithm
type RaftConfig struct {
	NodeID             string            // Unique identifier for this node
	PeerAddresses      map[string]string // Map of node IDs to network addresses
	ElectionTimeoutMin int               // Minimum election timeout in milliseconds
	ElectionTimeoutMax int               // Maximum election timeout in milliseconds
	HeartbeatInterval  int               // Heartbeat interval in milliseconds
}
