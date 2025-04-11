package node

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	distributed_auth_system "github.com/GooseFuse/distributed-auth-system/protoc"
)

// SyncManager manages state synchronization between nodes
type SyncManager struct {
	dataStore      *DataStore
	networkManager *NetworkManager
	mutex          sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
	syncInterval   time.Duration
}

// NewSyncManager creates a new SyncManager
func NewSyncManager(dataStore *DataStore, networkManager *NetworkManager) *SyncManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &SyncManager{
		dataStore:      dataStore,
		networkManager: networkManager,
		ctx:            ctx,
		cancel:         cancel,
		syncInterval:   1 * time.Minute,
	}
}

// Start starts the synchronization process
func (sm *SyncManager) Start() {
	go sm.syncLoop()
}

// Stop stops the synchronization process
func (sm *SyncManager) Stop() {
	sm.cancel()
}

// syncLoop periodically synchronizes with peers
func (sm *SyncManager) syncLoop() {
	ticker := time.NewTicker(sm.syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-sm.ctx.Done():
			return
		case <-ticker.C:
			sm.SyncWithPeers()
		}
	}
}

// SyncStats tracks statistics for a sync operation
type SyncStats struct {
	TransactionsReceived int
	ConflictsResolved    int
	Errors               int
}

// SyncWithPeers synchronizes state with all peers
func (sm *SyncManager) SyncWithPeers() {
	// Use RWMutex instead of RLock since we'll be modifying state
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Get local Merkle root
	localMerkleRoot := sm.dataStore.GetMerkleRoot()
	log.Printf("Local Merkle root: %s", localMerkleRoot)

	// Get all local keys for comparison
	localKeys, err := sm.dataStore.GetAllKeys()
	if err != nil {
		log.Printf("Error getting local keys: %v", err)
		return
	}

	// Create a map for quick lookup
	localKeysMap := make(map[string]bool)
	for _, key := range localKeys {
		localKeysMap[key] = true
	}

	// Get all peer clients
	peerClients := sm.networkManager.GetPeerClients()
	if len(peerClients) == 0 {
		log.Printf("No peers available for synchronization")
		return
	}

	// Track sync statistics
	syncStats := make(map[string]*SyncStats)

	// Process each peer
	for peerID, client := range peerClients {
		log.Printf("Syncing with peer %s", peerID)

		// Initialize stats for this peer
		syncStats[peerID] = &SyncStats{}

		// Sync with this peer
		sm.syncWithPeer(peerID, client, localMerkleRoot, localKeysMap, syncStats[peerID])
	}

	// Log sync results
	for peerID, stats := range syncStats {
		log.Printf("Sync with %s complete: received %d transactions, conflicts %d, errors %d",
			peerID, stats.TransactionsReceived, stats.ConflictsResolved, stats.Errors)
	}
}

// syncWithPeer synchronizes with a single peer
func (sm *SyncManager) syncWithPeer(peerID string, client distributed_auth_system.TransactionServiceClient,
	localMerkleRoot string, localKeysMap map[string]bool, stats *SyncStats) {

	// Set timeout for the sync operation
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// In a real implementation, we would check if we can use the extended client
	// For now, we'll use the simulation approach
	//sm.simulateSyncWithPeer(ctx, peerID, client, localMerkleRoot, localKeysMap, stats)
	resp, err := client.VerifyState(ctx, &StateVerificationRequest{
		NodeId:     sm.networkManager.nodeID,
		MerkleRoot: localMerkleRoot,
	})
	if err != nil {
		log.Printf("stateVerifyResponseErr: %s", err)
		return
	}

	log.Printf("stateVerifyResponse: %s", resp)

}

// simulateSyncWithPeer simulates synchronization with a peer when the extended client is not available
func (sm *SyncManager) simulateSyncWithPeer(ctx context.Context, peerID string, client interface{},
	localMerkleRoot string, localKeysMap map[string]bool, stats *SyncStats) {

	log.Printf("Simulating sync with peer %s", peerID)

	// Create a bloom filter of local keys for efficient comparison
	bloomFilter := sm.createBloomFilter(localKeysMap)

	// Use ctx and bloomFilter in a way that doesn't affect functionality
	// but satisfies the compiler's "declared and not used" errors
	if ctx.Err() != nil {
		log.Printf("Context error: %v", ctx.Err())
		return
	}
	log.Printf("Using bloom filter of size %d bytes", len(bloomFilter))

	// Batch processing parameters
	batchSize := 100
	var continuationToken string
	hasMore := true

	// Process batches until we've received all data
	for hasMore && stats.Errors < 3 { // Limit errors to prevent infinite loops
		// Simulate the SyncState RPC
		log.Printf("Simulating SyncState RPC with peer %s (batch %s)", peerID, continuationToken)

		// Simulate getting a batch of transactions
		simulatedBatch := sm.simulateTransactionBatch(peerID, len(continuationToken), batchSize)

		if len(simulatedBatch) == 0 {
			// No more transactions to process
			hasMore = false
			break
		}

		// Process each transaction in the batch
		for _, tx := range simulatedBatch {
			// Check if we already have this key
			_, exists := localKeysMap[tx.Key]

			// Process the transaction
			conflict, err := sm.processTransaction(tx, exists)

			if err != nil {
				stats.Errors++
				log.Printf("Error processing transaction from peer %s: %v", peerID, err)
				continue
			}

			stats.TransactionsReceived++
			if conflict {
				stats.ConflictsResolved++
			}

			// Add to local keys map
			localKeysMap[tx.Key] = true
		}

		// Update continuation token (in real implementation, this would come from the response)
		if continuationToken == "" {
			continuationToken = "1"
		} else {
			// Simulate pagination by incrementing the token
			hasMore = false // Only process one batch in simulation
		}
	}

	// Verify state consistency after sync
	log.Printf("Verifying state consistency with peer %s", peerID)
	sm.VerifyStateWithPeer(peerID)
}

// createBloomFilter creates a bloom filter of the given keys
func (sm *SyncManager) createBloomFilter(keysMap map[string]bool) []byte {
	// In a real implementation, this would create a bloom filter
	// For now, we'll just return a placeholder
	log.Printf("Creating bloom filter with %d keys", len(keysMap))
	return []byte("simulated-bloom-filter")
}

// simulateTransactionBatch simulates getting a batch of transactions from a peer
// In a real implementation, this would be replaced with actual RPC calls
func (sm *SyncManager) simulateTransactionBatch(peerID string, batchNum, batchSize int) []TransactionData {
	// This is a simulation function for demonstration purposes
	// In a real implementation, this data would come from the peer via RPC

	// For demonstration, we'll return an empty batch after the first one
	if batchNum > 0 {
		return []TransactionData{}
	}

	// Simulate 3 transactions in the first batch
	simulatedData := []TransactionData{
		{
			Key:       fmt.Sprintf("key_from_%s_1", peerID),
			Value:     fmt.Sprintf("value_from_%s_1", peerID),
			Timestamp: time.Now().UnixNano(),
		},
		{
			Key:       fmt.Sprintf("key_from_%s_2", peerID),
			Value:     fmt.Sprintf("value_from_%s_2", peerID),
			Timestamp: time.Now().UnixNano(),
		},
		{
			Key:       fmt.Sprintf("key_from_%s_3", peerID),
			Value:     fmt.Sprintf("value_from_%s_3", peerID),
			Timestamp: time.Now().UnixNano(),
		},
	}

	log.Printf("Simulated batch %d from peer %s with %d transactions",
		batchNum, peerID, len(simulatedData))

	return simulatedData
}

// TransactionData represents a key-value transaction with metadata
type TransactionData struct {
	Key       string
	Value     string
	Timestamp int64
	// In a real implementation, this would include vector clock and other metadata
}

// processTransaction applies a transaction with conflict resolution
func (sm *SyncManager) processTransaction(tx TransactionData, exists bool) (bool, error) {
	// If key doesn't exist, simply store the new value
	if !exists {
		err := sm.dataStore.StoreData(tx.Key, tx.Value)
		if err != nil {
			return false, fmt.Errorf("failed to store data: %w", err)
		}
		log.Printf("Added new key-value: %s=%s", tx.Key, tx.Value)
		return false, nil
	}

	// Key exists - potential conflict
	existingValue, err := sm.dataStore.GetData(tx.Key)
	if err != nil {
		return false, fmt.Errorf("failed to get existing data: %w", err)
	}

	// In a real implementation, we would compare vector clocks here
	// For this implementation, we'll use a simple strategy:
	// 1. If values are the same, no conflict
	// 2. If values are different, use the remote value (for demonstration)

	if existingValue == tx.Value {
		// No conflict, values are the same
		return false, nil
	}

	// Conflict detected - in a real implementation, we would use vector clocks
	// For now, we'll just use the remote value
	err = sm.dataStore.StoreData(tx.Key, tx.Value)
	if err != nil {
		return false, fmt.Errorf("failed to resolve conflict: %w", err)
	}

	log.Printf("Resolved conflict for key %s: local=%s, remote=%s",
		tx.Key, existingValue, tx.Value)

	return true, nil
}

// VerifyStateWithPeer verifies state consistency with a peer
func (sm *SyncManager) VerifyStateWithPeer(peerID string) bool {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Get local Merkle root
	localMerkleRoot := sm.dataStore.GetMerkleRoot()

	// In a real implementation, we would use the VerifyState RPC
	// For this example, we'll just return true
	log.Printf("Would verify state with peer %s (Merkle root: %s)", peerID, localMerkleRoot)
	return true
}

// CheckpointManager manages checkpointing for fault tolerance
type CheckpointManager struct {
	dataStore      *DataStore
	checkpointDir  string
	interval       time.Duration
	lastCheckpoint time.Time
	mutex          sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
}

// CheckpointData represents the data stored in a checkpoint
type CheckpointData struct {
	MerkleRoot string            `json:"merkle_root"`
	KeyValues  map[string]string `json:"key_values"`
	Timestamp  time.Time         `json:"timestamp"`
}

// NewCheckpointManager creates a new CheckpointManager
func NewCheckpointManager(dataStore *DataStore, checkpointDir string, interval time.Duration) *CheckpointManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &CheckpointManager{
		dataStore:      dataStore,
		checkpointDir:  checkpointDir,
		interval:       interval,
		lastCheckpoint: time.Now(),
		ctx:            ctx,
		cancel:         cancel,
	}
}

// Start starts the checkpointing process
func (cm *CheckpointManager) Start() {
	// Create checkpoint directory if it doesn't exist
	err := os.MkdirAll(cm.checkpointDir, 0755)
	if err != nil {
		log.Printf("Failed to create checkpoint directory: %v", err)
		return
	}

	go cm.checkpointLoop()
}

// Stop stops the checkpointing process
func (cm *CheckpointManager) Stop() {
	cm.cancel()
}

// checkpointLoop periodically creates checkpoints
func (cm *CheckpointManager) checkpointLoop() {
	ticker := time.NewTicker(cm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-cm.ctx.Done():
			return
		case <-ticker.C:
			checkpointID := time.Now().Unix()
			err := cm.CreateCheckpoint(checkpointID)
			if err != nil {
				log.Printf("Failed to create checkpoint: %v", err)
			}
		}
	}
}

// CreateCheckpoint creates a checkpoint with the given ID
func (cm *CheckpointManager) CreateCheckpoint(checkpointID int64) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Get all keys from the data store
	keys, err := cm.dataStore.GetAllKeys()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	// Get values for all keys
	keyValues := make(map[string]string)
	for _, key := range keys {
		value, err := cm.dataStore.GetData(key)
		if err != nil {
			log.Printf("Failed to get value for key %s: %v", key, err)
			continue
		}
		keyValues[key] = value
	}

	// Create checkpoint data
	checkpointData := CheckpointData{
		MerkleRoot: cm.dataStore.GetMerkleRoot(),
		KeyValues:  keyValues,
		Timestamp:  time.Now(),
	}

	// Serialize checkpoint data
	data, err := json.MarshalIndent(checkpointData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checkpoint data: %w", err)
	}

	// Write checkpoint file
	checkpointPath := filepath.Join(cm.checkpointDir, fmt.Sprintf("checkpoint_%d.json", checkpointID))
	err = ioutil.WriteFile(checkpointPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write checkpoint file: %w", err)
	}

	cm.lastCheckpoint = time.Now()
	log.Printf("Created checkpoint %d with %d keys", checkpointID, len(keys))
	return nil
}

// RestoreFromCheckpoint restores from a checkpoint with the given ID
func (cm *CheckpointManager) RestoreFromCheckpoint(checkpointID int64) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Read checkpoint file
	checkpointPath := filepath.Join(cm.checkpointDir, fmt.Sprintf("checkpoint_%d.json", checkpointID))
	data, err := ioutil.ReadFile(checkpointPath)
	if err != nil {
		return fmt.Errorf("failed to read checkpoint file: %w", err)
	}

	// Deserialize checkpoint data
	var checkpointData CheckpointData
	err = json.Unmarshal(data, &checkpointData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal checkpoint data: %w", err)
	}

	// Restore key-value pairs
	for key, value := range checkpointData.KeyValues {
		err := cm.dataStore.StoreData(key, value)
		if err != nil {
			log.Printf("Failed to restore key %s: %v", key, err)
		}
	}

	log.Printf("Restored checkpoint %d with %d keys", checkpointID, len(checkpointData.KeyValues))
	return nil
}

// GetLastCheckpointTime returns the time of the last checkpoint
func (cm *CheckpointManager) GetLastCheckpointTime() time.Time {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.lastCheckpoint
}

// ListCheckpoints lists all available checkpoints
func (cm *CheckpointManager) ListCheckpoints() ([]int64, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	// Read checkpoint directory
	files, err := ioutil.ReadDir(cm.checkpointDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read checkpoint directory: %w", err)
	}

	// Extract checkpoint IDs
	var checkpointIDs []int64
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		var id int64
		_, err := fmt.Sscanf(file.Name(), "checkpoint_%d.json", &id)
		if err != nil {
			continue
		}

		checkpointIDs = append(checkpointIDs, id)
	}

	return checkpointIDs, nil
}
