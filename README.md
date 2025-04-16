# Distributed Authorization and Authentication System

A distributed authorization and authentication system inspired by blockchain node architectures, emphasizing data consistency, network communication, and consensus management.

## Features

- **Data Storage**
  - Key-value data store with Merkle tree structure for efficient data integrity validation
  - LevelDB for lightweight, persistent storage
  - Redis caching layer for optimized read operations
  - Local ledger that periodically syncs with trusted peers

- **Network Communication**
  - gRPC for efficient, bidirectional communication between nodes
  - Pub/sub model with publisher and subscriber goroutines
  - Rate-limiting and backoff strategies for network stability

- **Consensus Mechanism**
  - Lightweight Raft consensus algorithm for data consistency
  - Leader elections and quorum maintenance
  - Majority consensus for transaction commitment

- **Synchronization and Fault Tolerance**
  - Checkpointing system for state snapshots
  - Hash-based state verification for data synchronization
  - Bloom filters for efficient identification of missing data

- **Security Enhancements**
  - End-to-end encryption using TLS
  - ECDSA signatures for transaction authentication

## Architecture

The system is composed of the following components:

1. **DataStore**: Manages data storage using LevelDB and Redis
2. **MerkleTree**: Provides data integrity validation
3. **ConsensusManager**: Implements the Raft consensus algorithm
4. **NetworkManager**: Handles communication between nodes
5. **SecurityManager**: Manages TLS and ECDSA signatures
6. **SyncManager**: Handles state synchronization between nodes
7. **CheckpointManager**: Manages checkpointing for fault tolerance

## Requirements

- Go 1.23 or higher
- Redis server
- LevelDB

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/GooseFuse/distributed-auth-system.git
   cd distributed-auth-system
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Generate gRPC code (if needed):
   ```
   protoc --go_out=. --go-grpc_out=. transaction.proto
   ```

## Running the System

### Starting a Node

```
go run . -node=node1 -port=:50051 -db=./data -redis=localhost:6379 -tls=false
```

Command-line flags:
- `-node`: Node ID (default: "node1")
- `-port`: Port to listen on (default: ":50051")
- `-db`: Path to database (default: "./data")
- `-redis`: Redis address (default: "localhost:6379")
- `-tls`: Use TLS (default: false)

### Running Multiple Nodes

To run a cluster of nodes, start multiple instances with different node IDs and ports:

```
# Terminal 1
go run . -node=node1 -port=:50051

# Terminal 2
go run . -node=node2 -port=:50052

# Terminal 3
go run . -node=node3 -port=:50053
```

### Using the Client

A client application is provided to interact with the system. The client allows you to store data, retrieve data, and simulate user authentication.

To build and run the client:

```
# Navigate to the client directory
cd client

# Download dependencies
go mod tidy

# Build the client
go build -o auth-client

# Run the client (examples)
./auth-client -op=store -key=user123 -value=password456
./auth-client -op=get -key=user123
./auth-client -op=auth -key=user123 -value=password456
```

For more details on using the client, see the [client README](client/README.md).

### Running the Demo

A demo script is provided to demonstrate how to run a cluster of nodes and interact with them using the client:

```
# Run the demo script
run_demo.bat
```

The demo script:
1. Builds the distributed auth system and client
2. Starts Redis if it's not already running
3. Creates a nodes directory with subdirectories for each node
4. Starts three nodes on different ports
5. Demonstrates client operations (store, get, authenticate)
6. Keeps the nodes running until you press any key
7. Automatically cleans up by stopping all node processes

## System Components

### Data Storage

The data storage layer uses LevelDB for persistent storage and Redis for caching. It also implements a Merkle tree for data integrity validation.

```go
// Store data
dataStore.StoreData("key", "value")

// Retrieve data
value, err := dataStore.GetData("key")
```

### Consensus

The consensus layer implements the Raft algorithm for leader election and log replication.

```go
// Propose a transaction
success, err := consensusManager.ProposeTransaction("key", "value")
```

### Network Communication

The network layer handles communication between nodes using gRPC.

```go
// Broadcast a transaction to all peers
networkManager.BroadcastTransaction("key", "value")
```

### Security

The security layer provides TLS encryption and ECDSA signatures.

```go
// Sign data
signature, err := securityManager.SignData(data)

// Verify signature
isValid := securityManager.VerifySignature(data, signature, publicKey)
```

## Development

### Project Structure

- `data_storage.go`: Data storage implementation
- `consensus.go`: Raft consensus implementation
- `network.go`: Network communication implementation
- `security.go`: Security features implementation
- `sync.go`: Synchronization and checkpointing implementation
- `grpc_service.go`: gRPC service implementation
- `main.go`: Main application entry point
- `transaction.proto`: Protocol buffer definition
- `client/`: Client application for interacting with the system
- `run_demo.bat`: Demo script for running a cluster of nodes
- `.gitignore`: Specifies files and directories to be excluded from version control
- `LICENSE`: MIT license file

### Adding New Features

1. Define the feature in the appropriate file
2. Update the gRPC service if needed
3. Add tests for the new feature
4. Update documentation

### Version Control

The project includes a `.gitignore` file that excludes the following from version control:

- Compiled binaries and executables
- Test files and coverage reports
- Dependency directories
- IDE-specific files
- OS-specific files
- Data directories and nodes directory (data/, nodes/)
- Database files (*.db, *.log, *.ldb)
- Redis dumps
- Certificate files
- Checkpoint files
- Temporary files
- Environment variable files

This ensures that only the source code and necessary configuration files are committed to the repository.


## License

This project is licensed under the MIT License - see the LICENSE file for details.



generate proto: `protoc --proto_path=. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. protoc/transaction.proto`
