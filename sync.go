package main

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

// SyncWithPeers synchronizes state with all peers
func (sm *SyncManager) SyncWithPeers() {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Get local Merkle root
	localMerkleRoot := sm.dataStore.GetMerkleRoot()

	// Sync with each peer
	peerClients := sm.networkManager.GetPeerClients()
	for peerID := range peerClients {
		log.Printf("Syncing with peer %s", peerID)

		// In a real implementation, we would use the SyncState RPC
		// For this example, we'll just log the sync attempt
		log.Printf("Would sync with peer %s (Merkle root: %s)", peerID, localMerkleRoot)
	}
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
