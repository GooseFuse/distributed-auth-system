package datastore

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/GooseFuse/distributed-auth-system/protoc"
	"github.com/go-redis/redis/v8"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// MerkleTree represents a Merkle tree structure for data integrity validation
type MerkleTree struct {
	Root      *MerkleNode
	mutex     sync.RWMutex
	leafNodes map[string]*MerkleNode // Map of data hash to leaf node for quick lookups
}

// NewMerkleTree creates a new Merkle tree
func NewMerkleTree() *MerkleTree {
	return &MerkleTree{
		Root:      nil,
		leafNodes: make(map[string]*MerkleNode),
	}
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Hash  string
	Left  *MerkleNode
	Right *MerkleNode
	Data  string // Only leaf nodes have data
}

// NewMerkleNode creates a new Merkle node
func NewMerkleNode(left, right *MerkleNode, data string) *MerkleNode {
	node := &MerkleNode{
		Left:  left,
		Right: right,
		Data:  data,
	}

	// If it's a leaf node (no children), hash the data
	if left == nil && right == nil {
		node.Hash = Hash(data)
	} else {
		// Otherwise, hash the concatenation of the children's hashes
		prevHashes := left.Hash + right.Hash
		node.Hash = Hash(prevHashes)
	}

	return node
}

// AddLeaf adds a new leaf node to the Merkle tree
func (mt *MerkleTree) AddLeaf(data string) error {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()

	dataHash := Hash(data)

	// Check if the leaf already exists
	if _, exists := mt.leafNodes[dataHash]; exists {
		return nil // Leaf already exists, no need to add it again
	}

	// Create a new leaf node
	newNode := NewMerkleNode(nil, nil, data)
	mt.leafNodes[dataHash] = newNode

	// If this is the first node, make it the root
	if mt.Root == nil {
		mt.Root = newNode
		return nil
	}

	// Otherwise, rebuild the tree with the new leaf
	mt.rebuildTree()
	return nil
}

// rebuildTree rebuilds the Merkle tree from the leaf nodes
func (mt *MerkleTree) rebuildTree() {
	nodes := make([]*MerkleNode, 0, len(mt.leafNodes))

	// Collect all leaf nodes
	for _, node := range mt.leafNodes {
		nodes = append(nodes, node)
	}

	// Build the tree bottom-up
	for len(nodes) > 1 {
		var newLevel []*MerkleNode

		// Process nodes two at a time
		for i := 0; i < len(nodes); i += 2 {
			// If we have an odd number of nodes, duplicate the last one
			if i+1 == len(nodes) {
				newNode := NewMerkleNode(nodes[i], nodes[i], "")
				newLevel = append(newLevel, newNode)
			} else {
				newNode := NewMerkleNode(nodes[i], nodes[i+1], "")
				newLevel = append(newLevel, newNode)
			}
		}

		nodes = newLevel
	}

	// The last remaining node is the root
	if len(nodes) > 0 {
		mt.Root = nodes[0]
	}
}

// VerifyData verifies if data exists in the Merkle tree
func (mt *MerkleTree) VerifyData(data string) bool {
	mt.mutex.RLock()
	defer mt.mutex.RUnlock()

	dataHash := Hash(data)
	_, exists := mt.leafNodes[dataHash]
	return exists
}

// GetRoot returns the root hash of the Merkle tree
func (mt *MerkleTree) GetRoot() string {
	mt.mutex.RLock()
	defer mt.mutex.RUnlock()

	if mt.Root == nil {
		return ""
	}
	return mt.Root.Hash
}

// BloomFilter represents a Bloom filter for efficient data lookups
type BloomFilter struct {
	filter        []byte
	hashFunctions int
}

// NewBloomFilter creates a new Bloom filter
func NewBloomFilter(size int, hashFunctions int) *BloomFilter {
	return &BloomFilter{
		filter:        make([]byte, size),
		hashFunctions: hashFunctions,
	}
}

// Add adds an item to the Bloom filter
func (bf *BloomFilter) Add(data string) {
	for i := 0; i < bf.hashFunctions; i++ {
		// Use different hash seeds for each hash function
		hashValue := Hash(fmt.Sprintf("%s:%d", data, i))
		// Convert first 4 bytes of hash to uint32 and mod by filter size
		index := uint32(hashValue[0])<<24 | uint32(hashValue[1])<<16 | uint32(hashValue[2])<<8 | uint32(hashValue[3])
		index %= uint32(len(bf.filter))
		bf.filter[index] = 1
	}
}

// Contains checks if an item might be in the Bloom filter
func (bf *BloomFilter) Contains(data string) bool {
	for i := 0; i < bf.hashFunctions; i++ {
		hashValue := Hash(fmt.Sprintf("%s:%d", data, i))
		index := uint32(hashValue[0])<<24 | uint32(hashValue[1])<<16 | uint32(hashValue[2])<<8 | uint32(hashValue[3])
		index %= uint32(len(bf.filter))
		if bf.filter[index] == 0 {
			return false
		}
	}
	return true
}

// DataStore represents the key-value store with LevelDB and Redis
type DataStore struct {
	db          *leveldb.DB
	cache       *redis.Client
	merkleTree  *MerkleTree
	bloomFilter *BloomFilter
	ctx         context.Context
	mutex       sync.RWMutex
	dbPath      string
}

func (d *DataStore) GetDBPath() string {
	return d.dbPath
}

// NewDataStore initializes a new DataStore
func NewDataStore(dbPath string, redisAddr string) (*DataStore, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}

	cache := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx := context.Background()

	// Test Redis connection
	_, err = cache.Ping(ctx).Result()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	merkleTree := NewMerkleTree()
	bloomFilter := NewBloomFilter(10000, 5) // 10KB filter with 5 hash functions

	ds := &DataStore{
		db:          db,
		cache:       cache,
		merkleTree:  merkleTree,
		bloomFilter: bloomFilter,
		ctx:         ctx,
		dbPath:      dbPath,
	}

	// Initialize Merkle tree and Bloom filter with existing data
	err = ds.initializeFromDB()
	if err != nil {
		db.Close()
		cache.Close()
		return nil, err
	}

	return ds, nil
}

// initializeFromDB loads existing data from LevelDB into the Merkle tree and Bloom filter
func (ds *DataStore) initializeFromDB() error {
	iter := ds.db.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())

		// Add to Merkle tree
		err := ds.merkleTree.AddLeaf(key + ":" + value)
		if err != nil {
			return err
		}

		// Add to Bloom filter
		ds.bloomFilter.Add(key)
	}

	return iter.Error()
}

// StoreData stores data in LevelDB and Redis cache, and updates the Merkle tree
func (ds *DataStore) StoreData(key string, value string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	// Store in LevelDB
	err := ds.db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		return err
	}

	// Store in Redis cache with expiration
	err = ds.cache.Set(ds.ctx, key, value, 1*time.Hour).Err()
	if err != nil {
		// Log the error but continue, as cache failures are not critical
		fmt.Printf("Redis cache error: %v\n", err)
	}

	// Update Merkle tree
	err = ds.merkleTree.AddLeaf(key + ":" + value)
	if err != nil {
		return err
	}

	// Update Bloom filter
	ds.bloomFilter.Add(key)

	return nil
}

// GetData retrieves data from Redis cache or LevelDB
func (ds *DataStore) GetData(key string) (string, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	// Check Bloom filter first for quick negative lookups
	if !ds.bloomFilter.Contains(key) {
		return "", errors.New("key not found")
	}

	// Try to get from Redis cache first
	val, err := ds.cache.Get(ds.ctx, key).Result()
	if err == nil {
		return val, nil
	}

	// If not in cache or error occurred, get from LevelDB
	data, err := ds.db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}

	// Update cache for future reads
	ds.cache.Set(ds.ctx, key, string(data), 1*time.Hour)

	return string(data), nil
}

// GetMerkleRoot returns the root hash of the Merkle tree
func (ds *DataStore) GetMerkleRoot() string {
	return ds.merkleTree.GetRoot()
}

// VerifyData verifies if a key-value pair exists in the store
func (ds *DataStore) VerifyData(key string, value string) bool {
	return ds.merkleTree.VerifyData(key + ":" + value)
}

// GetAllKeys returns all keys in the database
func (ds *DataStore) GetAllKeys() ([]string, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	var keys []string
	iter := ds.db.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		keys = append(keys, string(iter.Key()))
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}

	return keys, nil
}

// GetKeysByPrefix returns all keys with a given prefix
func (ds *DataStore) GetKeysByPrefix(prefix string) ([]string, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	var keys []string
	iter := ds.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	for iter.Next() {
		keys = append(keys, string(iter.Key()))
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (ds *DataStore) GetAllTransactions() ([]*protoc.Transaction, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	var txs []*protoc.Transaction
	iter := ds.db.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		txs = append(txs, &protoc.Transaction{Key: string(iter.Key()), Value: string(iter.Value())})
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}

	return txs, nil
}

// Close closes the database connections
func (ds *DataStore) Close() {
	ds.db.Close()
	ds.cache.Close()
}

// Hash generates a SHA-256 hash of the input data
func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
