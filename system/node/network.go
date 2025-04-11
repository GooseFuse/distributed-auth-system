package node

import (
	"context"
	"log"
	"sync"
	"time"

	distributed_auth_system "github.com/GooseFuse/distributed-auth-system/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// NetworkManager manages network communication between nodes
type NetworkManager struct {
	nodeID          string
	peerNodes       map[string]string // Map of node IDs to addresses
	peerClients     map[string]distributed_auth_system.TransactionServiceClient
	pubSubManager   *PubSubManager
	rateLimiter     *RateLimiter
	securityManager *SecurityManager
	mutex           sync.RWMutex
	ctx             context.Context
	cancel          context.CancelFunc
}

// NewNetworkManager creates a new NetworkManager
func NewNetworkManager(nodeID string, securityManager *SecurityManager) *NetworkManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &NetworkManager{
		nodeID:          nodeID,
		peerNodes:       make(map[string]string),
		peerClients:     make(map[string]distributed_auth_system.TransactionServiceClient),
		pubSubManager:   NewPubSubManager(),
		rateLimiter:     NewRateLimiter(100, 1*time.Minute), // 100 requests per minute
		securityManager: securityManager,
		ctx:             ctx,
		cancel:          cancel,
	}
}

// Start starts the network manager
func (nm *NetworkManager) Start() {
	// Connect to all peers
	for peerID, peerAddr := range nm.peerNodes {
		if peerID != nm.nodeID {
			go nm.ConnectToPeer(peerID, peerAddr)
		}
	}
}

// Stop stops the network manager
func (nm *NetworkManager) Stop() {
	nm.cancel()
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	// Close all connections
	// In a real implementation, we would close the gRPC connections
	log.Printf("Stopping network manager, disconnecting from %d peers", len(nm.peerClients))
}

// AddPeer adds a peer to the network
func (nm *NetworkManager) AddPeer(peerID string, peerAddr string) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	nm.peerNodes[peerID] = peerAddr
	go nm.ConnectToPeer(peerID, peerAddr)
}

// RemovePeer removes a peer from the network
func (nm *NetworkManager) RemovePeer(peerID string) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	delete(nm.peerNodes, peerID)
	delete(nm.peerClients, peerID)
}

// ConnectToPeer connects to a peer
func (nm *NetworkManager) ConnectToPeer(peerID string, peerAddr string) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	// Check if already connected
	if _, ok := nm.peerClients[peerID]; ok {
		return
	}

	// Set up connection options
	var opts []grpc.DialOption
	if nm.securityManager != nil {
		// Use TLS
		creds := credentials.NewTLS(nm.securityManager.GetTLSConfig())
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		// Use insecure connection (for development only)
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient("peerAddr", opts...)
	if err != nil {
		log.Printf("Failed to connect to peer %s at %s: %v", peerID, peerAddr, err)
		return
	}

	// Create client
	nm.peerClients[peerID] = distributed_auth_system.NewTransactionServiceClient(conn)

	log.Printf("Connected to peer %s at %s", peerID, peerAddr)
}

// GetPeerClients returns a map of peer clients
func (nm *NetworkManager) GetPeerClients() map[string]distributed_auth_system.TransactionServiceClient {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	// Create a copy of the map to avoid concurrent access issues
	clients := make(map[string]distributed_auth_system.TransactionServiceClient)
	for id, client := range nm.peerClients {
		clients[id] = client
	}

	return clients
}

// BroadcastTransaction broadcasts a key-value transaction to all peers
func (nm *NetworkManager) BroadcastTransaction(key, value string) {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	// Create transaction request
	req := &distributed_auth_system.TransactionRequest{
		Key:   key,
		Value: value,
	}

	// Broadcast to all peers
	for peerID, client := range nm.peerClients {
		go func(id string, c distributed_auth_system.TransactionServiceClient) {
			ctx, cancel := context.WithTimeout(nm.ctx, 5*time.Second)
			defer cancel()

			_, err := c.HandleTransaction(ctx, req)
			if err != nil {
				log.Printf("Failed to broadcast transaction to peer %s: %v", id, err)
			}
		}(peerID, client)
	}

	// Publish to subscribers
	nm.pubSubManager.Publish("transactions", map[string]string{
		"key":   key,
		"value": value,
	})
}

// PubSubManager manages the publish-subscribe messaging system
type PubSubManager struct {
	publishers  map[string]chan interface{}
	subscribers map[string][]chan interface{}
	mutex       sync.RWMutex
}

// NewPubSubManager creates a new PubSubManager
func NewPubSubManager() *PubSubManager {
	return &PubSubManager{
		publishers:  make(map[string]chan interface{}),
		subscribers: make(map[string][]chan interface{}),
	}
}

// Publish publishes a message to a topic
func (pm *PubSubManager) Publish(topic string, message interface{}) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	// Get publisher channel
	publisher, ok := pm.publishers[topic]
	if !ok {
		// Create publisher channel if it doesn't exist
		publisher = make(chan interface{}, 100)
		pm.publishers[topic] = publisher

		// Start publisher goroutine
		go func() {
			for msg := range publisher {
				pm.mutex.RLock()
				subscribers := pm.subscribers[topic]
				pm.mutex.RUnlock()

				for _, subscriber := range subscribers {
					select {
					case subscriber <- msg:
						// Message sent
					default:
						// Subscriber channel is full, drop message
					}
				}
			}
		}()
	}

	// Send message to publisher channel
	select {
	case publisher <- message:
		// Message sent
	default:
		// Publisher channel is full, drop message
		log.Printf("Publisher channel for topic %s is full, dropping message", topic)
	}
}

// Subscribe subscribes to a topic
func (pm *PubSubManager) Subscribe(topic string) chan interface{} {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Create subscriber channel
	subscriber := make(chan interface{}, 100)

	// Add subscriber to topic
	pm.subscribers[topic] = append(pm.subscribers[topic], subscriber)

	return subscriber
}

// Unsubscribe unsubscribes from a topic
func (pm *PubSubManager) Unsubscribe(topic string, subscriber chan interface{}) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Get subscribers for topic
	subscribers := pm.subscribers[topic]
	if subscribers == nil {
		return
	}

	// Remove subscriber from topic
	for i, sub := range subscribers {
		if sub == subscriber {
			pm.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
			close(subscriber)
			break
		}
	}
}

// RateLimiter manages rate limiting for network requests
type RateLimiter struct {
	limits     map[string]int
	windows    map[string]time.Time
	windowSize time.Duration
	maxLimit   int
	mutex      sync.RWMutex
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(maxLimit int, windowSize time.Duration) *RateLimiter {
	return &RateLimiter{
		limits:     make(map[string]int),
		windows:    make(map[string]time.Time),
		windowSize: windowSize,
		maxLimit:   maxLimit,
	}
}

// AllowRequest checks if a request from the given client is allowed based on rate limits
func (rl *RateLimiter) AllowRequest(clientID string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Check if the window has expired
	lastRequest, ok := rl.windows[clientID]
	if !ok || now.Sub(lastRequest) > rl.windowSize {
		// Reset the window
		rl.windows[clientID] = now
		rl.limits[clientID] = 1
		return true
	}

	// Check if the limit has been reached
	if rl.limits[clientID] >= rl.maxLimit {
		return false
	}

	// Increment the counter
	rl.limits[clientID]++
	return true
}

// SetLimit sets the rate limit for a client
func (rl *RateLimiter) SetLimit(clientID string, limit int) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	rl.limits[clientID] = limit
}

// GetLimit gets the rate limit for a client
func (rl *RateLimiter) GetLimit(clientID string) int {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	limit, ok := rl.limits[clientID]
	if !ok {
		return rl.maxLimit
	}
	return limit
}
