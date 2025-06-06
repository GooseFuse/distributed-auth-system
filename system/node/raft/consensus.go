package raft

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/GooseFuse/distributed-auth-system/protoc"
	"github.com/GooseFuse/distributed-auth-system/system/datastore"
	"github.com/GooseFuse/distributed-auth-system/system/interfaces"
)

// RaftState represents the state of a node in the Raft consensus algorithm
type RaftState int

const (
	// Follower state - receives log entries from leader
	Follower RaftState = iota
	// Candidate state - requests votes from other nodes
	Candidate
	// Leader state - coordinates log replication
	Leader
)

// RaftNode represents a node in the Raft consensus algorithm
type RaftNode struct {
	config RaftConfig

	// Persistent state on all servers
	currentTerm int                // Latest term server has seen
	votedFor    string             // CandidateID that received vote in current term (or empty if none)
	log         []*protoc.LogEntry // Log entries

	// Volatile state on all servers
	commitIndex int // Index of highest log entry known to be committed
	lastApplied int // Index of highest log entry applied to state machine

	// Volatile state on leaders
	nextIndex  map[string]int // For each server, index of the next log entry to send
	matchIndex map[string]int // For each server, index of highest log entry known to be replicated

	// Runtime state
	state              RaftState
	electionTimeout    time.Duration
	lastHeartbeat      time.Time
	leaderID           string
	applyCh            chan *protoc.LogEntry
	stopCh             chan struct{}
	dataStore          *datastore.DataStore
	transactionHandler func(*protoc.Transaction) error
	networkManager     interfaces.NetworkManagerI
	mutex              sync.RWMutex
}

// NewRaftNode creates a new Raft node
func NewRaftNode(config RaftConfig, dataStore *datastore.DataStore, transactionHandler func(*protoc.Transaction) error) *RaftNode {
	node := &RaftNode{
		config:             config,
		currentTerm:        0,
		votedFor:           "",
		log:                make([]*protoc.LogEntry, 0),
		commitIndex:        -1,
		lastApplied:        -1,
		nextIndex:          make(map[string]int),
		matchIndex:         make(map[string]int),
		state:              Follower,
		electionTimeout:    randomTimeout(config.ElectionTimeoutMin, config.ElectionTimeoutMax),
		lastHeartbeat:      time.Now(),
		leaderID:           "",
		applyCh:            make(chan *protoc.LogEntry, 100),
		stopCh:             make(chan struct{}),
		dataStore:          dataStore,
		transactionHandler: transactionHandler,
	}

	return node
}

// Start starts the Raft node
func (rn *RaftNode) Start() {
	go rn.run()
	go rn.applyCommittedEntries()
}

// Stop stops the Raft node
func (rn *RaftNode) Stop() {
	close(rn.stopCh)
}

// run is the main loop for the Raft node
func (rn *RaftNode) run() {
	for {
		select {
		case <-rn.stopCh:
			return
		default:
			rn.mutex.RLock()
			state := rn.state
			elapsed := time.Since(rn.lastHeartbeat)
			timeout := rn.electionTimeout
			rn.mutex.RUnlock()

			switch state {
			case Follower, Candidate:
				if elapsed >= timeout {
					rn.startElection()
				}
			case Leader:
				time.Sleep(time.Duration(rn.config.HeartbeatInterval) * time.Millisecond)
				rn.sendHeartbeats()
			}

			// Small sleep to prevent CPU spinning
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// applyCommittedEntries applies committed log entries to the state machine
func (rn *RaftNode) applyCommittedEntries() {
	for {
		select {
		case <-rn.stopCh:
			return
		case entry := <-rn.applyCh:
			// Apply the command to the state machine
			err := rn.transactionHandler(entry.Command)
			if err != nil {
				fmt.Printf("Error applying command: %v\n", err)
			}
		}
	}
}

// startElection starts a leader election
func (rn *RaftNode) startElection() {
	rn.mutex.Lock()
	rn.state = Candidate
	rn.currentTerm++
	rn.votedFor = rn.config.NodeID // Vote for self
	currentTerm := rn.currentTerm
	lastLogIndex := len(rn.log) - 1
	lastLogTerm := 0
	if lastLogIndex >= 0 {
		lastLogTerm = int(rn.log[lastLogIndex].Term)
	}
	rn.electionTimeout = randomTimeout(rn.config.ElectionTimeoutMin, rn.config.ElectionTimeoutMax)
	rn.lastHeartbeat = time.Now()
	defer rn.mutex.Unlock()

	peerClients := rn.networkManager.GetPeerClients()
	fmt.Printf("Node %s starting election for term %d with %d clients\n", rn.config.NodeID, currentTerm, len(peerClients))

	// Request votes from all peers
	votesReceived := 1 // Vote for self
	votesNeeded := (len(rn.config.PeerAddresses) / 2) + 1

	var votesMutex sync.Mutex
	var wg sync.WaitGroup

	for peerID, client := range peerClients {
		if peerID == rn.config.NodeID {
			log.Printf("🔁 Skipping self (%s) in peer loop", peerID)
			continue
		}
		log.Printf("🔁 PeerID: %s selfID: %s", peerID, rn.config.NodeID)
		wg.Add(1)
		go func(peerID string, client protoc.TransactionServiceClient) {
			defer wg.Done()

			granted, err := rn.RequestVote(peerID, client, currentTerm, lastLogIndex, lastLogTerm)
			if err != nil {
				fmt.Printf("❌ RequestVote RPC Error from %s: %v\n", peerID, err)
				return
			}

			if granted {
				votesMutex.Lock()
				votesReceived++
				votesMutex.Unlock()
			}
		}(peerID, client)
	}

	// Wait for all vote requests to complete
	wg.Wait()

	// Check if we won the election
	//rn.mutex.Lock()
	//defer rn.mutex.Unlock()

	// Only become leader if still a candidate and term hasn't changed
	if rn.state == Candidate && rn.currentTerm == currentTerm && votesReceived >= votesNeeded {
		fmt.Printf("Node %s is the new leader, won election for term %d with %d votes\n", rn.config.NodeID, currentTerm, votesReceived)
		rn.becomeLeader()
	}
}

// becomeLeader transitions the node to the leader state
func (rn *RaftNode) becomeLeader() {
	log.Printf("Become Leader...\n")
	rn.state = Leader
	rn.leaderID = rn.config.NodeID

	// Initialize nextIndex and matchIndex for all peers
	lastLogIndex := len(rn.log)
	peerClients := rn.networkManager.GetPeerClients()
	for peerID := range peerClients {
		rn.nextIndex[peerID] = lastLogIndex
		rn.matchIndex[peerID] = -1
	}

	// Send initial empty AppendEntries RPCs (heartbeats) to each server
	go rn.sendHeartbeats()
}

// sendHeartbeats sends heartbeats to all peers
func (rn *RaftNode) sendHeartbeats() {
	rn.mutex.RLock()
	if rn.state != Leader {
		rn.mutex.RUnlock()
		return
	}

	currentTerm := rn.currentTerm
	commitIndex := rn.commitIndex
	rn.mutex.RUnlock()

	var wg sync.WaitGroup

	peerClients := rn.networkManager.GetPeerClients()
	for peerID, client := range peerClients {
		if peerID == rn.leaderID {
			continue
		}
		wg.Add(1)
		go func(peerID string, client protoc.TransactionServiceClient) {
			defer wg.Done()

			rn.mutex.RLock()
			nextIdx := rn.nextIndex[peerID]
			prevLogIndex := nextIdx - 1
			prevLogTerm := 0
			if prevLogIndex >= 0 && prevLogIndex < len(rn.log) {
				prevLogTerm = int(rn.log[prevLogIndex].Term)
			}

			// Get log entries to send
			entries := make([]*protoc.LogEntry, 0)
			if nextIdx < len(rn.log) {
				entries = rn.log[nextIdx:]
			}
			rn.mutex.RUnlock()

			// Send AppendEntries RPC
			success, err := rn.appendEntries(client, currentTerm, prevLogIndex, prevLogTerm, entries, commitIndex)
			if err != nil {
				fmt.Printf("Error sending AppendEntries to %s: %v\n", peerID, err)
				return
			}
			if success {
				rn.mutex.Lock()
				if len(entries) > 0 {
					rn.nextIndex[peerID] = nextIdx + len(entries)
					rn.matchIndex[peerID] = rn.nextIndex[peerID] - 1
				}
				rn.mutex.Unlock()
			} else {
				// If AppendEntries fails because of log inconsistency, decrement nextIndex and retry
				rn.mutex.Lock()
				if rn.nextIndex[peerID] > 0 {
					rn.nextIndex[peerID]--
				}
				rn.mutex.Unlock()
			}
		}(peerID, client)
	}
	wg.Wait()
	// Update commit index if needed
	rn.updateCommitIndex()
}

// updateCommitIndex updates the commit index based on the matchIndex values
func (rn *RaftNode) updateCommitIndex() {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if rn.state != Leader {
		return
	}

	// Find the highest index that is replicated on a majority of servers
	for n := rn.commitIndex + 1; n < len(rn.log); n++ {
		if rn.log[n].Term != int32(rn.currentTerm) {
			continue
		}

		count := 1 // Leader has the entry
		for _, matchIdx := range rn.matchIndex {
			if matchIdx >= n {
				count++
			}
		}

		if count > (len(rn.config.PeerAddresses) / 2) {
			rn.commitIndex = n
			// Signal that there are new entries to apply
			for i := rn.lastApplied + 1; i <= rn.commitIndex; i++ {
				rn.applyCh <- rn.log[i]
				rn.lastApplied = i
			}
		}
	}
}

// ProposeCommand proposes a new command to the Raft cluster
func (rn *RaftNode) ProposeCommand(command *protoc.Transaction) (bool, error) {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if rn.state != Leader {
		return false, errors.New("not the leader")
	}

	// Append to local log
	entry := &protoc.LogEntry{
		Term:    int32(rn.currentTerm),
		Index:   int32(len(rn.log)),
		Command: command,
	}
	rn.log = append(rn.log, entry)

	fmt.Printf("Leader %s proposing command at index %d\n", rn.config.NodeID, entry.Index)

	// Trigger immediate replication
	go rn.sendHeartbeats()

	return true, nil
}

// requestVote sends a RequestVote RPC to a peer
func (rn *RaftNode) RequestVote(peerID string, client protoc.TransactionServiceClient, term, lastLogIndex, lastLogTerm int) (bool, error) {
	fmt.Printf("📩 [RequestVote] From: %s | Term: %d | LastLogIndex: %d | LastLogTerm: %d\n",
		peerID, term, lastLogIndex, lastLogTerm)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req := &protoc.RequestVoteRequest{
		Term:         int32(term),
		LastLogIndex: int32(lastLogIndex),
		LastLogTerm:  int32(lastLogTerm),
		CandidateId:  rn.config.NodeID,
	}

	resp, err := client.RequestVote(ctx, req)
	if err != nil {
		return false, err
	}

	rn.updateTerm(resp.Term)

	return resp.VoteGranted, nil
}

// appendEntries sends an AppendEntries RPC to a peer
func (rn *RaftNode) appendEntries(client protoc.TransactionServiceClient, term, prevLogIndex, prevLogTerm int, entries []*protoc.LogEntry, leaderCommit int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req := &protoc.AppendEntriesRequest{
		Term:         int32(term),
		LeaderId:     rn.leaderID,
		PrevLogIndex: int32(prevLogIndex),
		PrevLogTerm:  int32(prevLogTerm),
		Entries:      entries,
	}

	resp, err := client.AppendEntries(ctx, req)
	if err != nil {
		return false, err
	}

	rn.updateTerm(resp.Term)

	return resp.Success, nil
}

func (rn *RaftNode) updateTerm(term int32) {
	if term > int32(rn.currentTerm) {
		rn.currentTerm = int(term)
		rn.state = Follower
		rn.votedFor = ""
		log.Printf("🔄 Stepping down: higher term seen from peer (%d > %d)", term, rn.currentTerm)
	}
}

// HandleRequestVote handles a RequestVote RPC from a peer
func (rn *RaftNode) HanldeRequestVote(candidateID string, term, lastLogIndex, lastLogTerm int) (int, bool) {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	log.Printf("🔁 RequestVote request from candidateID %s | term: %d | lastLogIndex: %d | lastLogTerm: %d",
		candidateID, term, lastLogIndex, lastLogTerm)

	// If term > currentTerm, update currentTerm
	if term > rn.currentTerm {
		rn.currentTerm = term
		rn.state = Follower
		rn.votedFor = ""
	}

	// Reply false if term < currentTerm
	if term < rn.currentTerm {
		return rn.currentTerm, false
	}

	// Check if candidate's log is at least as up-to-date as receiver's log
	lastIndex := len(rn.log) - 1
	lastTerm := 0
	if lastIndex >= 0 {
		lastTerm = int(rn.log[lastIndex].Term)
	}

	logOk := (lastLogTerm > lastTerm) || (lastLogTerm == lastTerm && lastLogIndex >= lastIndex)

	// Grant vote if votedFor is null or candidateID, and candidate's log is at least as up-to-date
	if (rn.votedFor == "" || rn.votedFor == candidateID) && logOk {
		rn.votedFor = candidateID
		rn.lastHeartbeat = time.Now() // Reset election timeout
		return rn.currentTerm, true
	}

	return rn.currentTerm, false
}

// HandleAppendEntries handles an AppendEntries RPC from a peer
func (rn *RaftNode) HandleAppendEntries(leaderID string, term, prevLogIndex, prevLogTerm int, entries []*protoc.LogEntry, leaderCommit int) (int, bool) {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	// If term > currentTerm, update currentTerm
	if term > rn.currentTerm {
		rn.currentTerm = term
		rn.state = Follower
		rn.votedFor = ""
	}

	// Reply false if term < currentTerm
	if term < rn.currentTerm {
		return rn.currentTerm, false
	}

	// Update leader ID and reset election timeout
	rn.leaderID = leaderID
	rn.lastHeartbeat = time.Now()

	// If we were a candidate, step down
	if rn.state == Candidate {
		rn.state = Follower
	}

	// Reply false if log doesn't contain an entry at prevLogIndex whose term matches prevLogTerm
	if prevLogIndex >= 0 {
		if prevLogIndex >= len(rn.log) {
			return rn.currentTerm, false
		}
		if rn.log[prevLogIndex].Term != int32(prevLogTerm) {
			return rn.currentTerm, false
		}
	}

	// If an existing entry conflicts with a new one (same index but different terms),
	// delete the existing entry and all that follow it
	for i, entry := range entries {
		if prevLogIndex+1+i < len(rn.log) {
			if rn.log[prevLogIndex+1+i].Term != entry.Term {
				rn.log = rn.log[:prevLogIndex+1+i]
				break
			}
		} else {
			break
		}
	}

	// Append any new entries not already in the log
	for i, entry := range entries {
		if prevLogIndex+1+i >= len(rn.log) {
			rn.log = append(rn.log, entry)
		}
	}

	// Update commit index if leader commit > commitIndex
	if leaderCommit > rn.commitIndex {
		rn.commitIndex = min(leaderCommit, len(rn.log)-1)

		// Apply newly committed entries
		for i := rn.lastApplied + 1; i <= rn.commitIndex; i++ {
			rn.applyCh <- rn.log[i]
			rn.lastApplied = i
		}
	}

	return rn.currentTerm, true
}

// IsLeader returns true if this node is the leader
func (rn *RaftNode) IsLeader() bool {
	rn.mutex.RLock()
	defer rn.mutex.RUnlock()
	return rn.state == Leader
}

// GetLeader returns the ID of the current leader
func (rn *RaftNode) GetLeader() string {
	rn.mutex.RLock()
	defer rn.mutex.RUnlock()
	return rn.leaderID
}

// GetState returns the current state of the Raft node
func (rn *RaftNode) GetState() (int, bool) {
	rn.mutex.RLock()
	defer rn.mutex.RUnlock()
	return rn.currentTerm, rn.state == Leader
}

// randomTimeout generates a random timeout between min and max milliseconds
func randomTimeout(min, max int) time.Duration {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Duration(min+rng.Intn(max-min)) * time.Millisecond
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ConsensusManager manages the Raft consensus algorithm
type ConsensusManager struct {
	raftNode  *RaftNode
	dataStore *datastore.DataStore
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewConsensusManager creates a new ConsensusManager
func NewConsensusManager(config RaftConfig, dataStore *datastore.DataStore, networkManager interfaces.NetworkManagerI) *ConsensusManager {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a transaction handler function
	transactionHandler := func(command *protoc.Transaction) error {
		return dataStore.StoreData(command.Key, command.Value)
	}

	raftNode := NewRaftNode(config, dataStore, transactionHandler)
	// Get all peer clients
	raftNode.networkManager = networkManager
	peerClients := networkManager.GetPeerClients()
	if len(peerClients) == 0 {
		log.Printf("No peers available for consensus")
	}

	return &ConsensusManager{
		raftNode:  raftNode,
		dataStore: dataStore,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start starts the ConsensusManager
func (cm *ConsensusManager) Start() {
	cm.raftNode.Start()
}

// Stop stops the ConsensusManager
func (cm *ConsensusManager) Stop() {
	cm.raftNode.Stop()
	cm.cancel()
}

// ProposeTransaction proposes a transaction to the Raft cluster
func (cm *ConsensusManager) ProposeTransaction(key, value string) (bool, error) {
	return cm.raftNode.ProposeCommand(&protoc.Transaction{Key: key, Value: value})
}

// IsLeader returns true if this node is the leader
func (cm *ConsensusManager) IsLeader() bool {
	return cm.raftNode.IsLeader()
}

// GetLeader returns the ID of the current leader
func (cm *ConsensusManager) GetLeader() string {
	return cm.raftNode.GetLeader()
}
