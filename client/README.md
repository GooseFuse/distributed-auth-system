# Distributed Auth System Client

This is a simple client for interacting with the Distributed Authorization and Authentication System.

## Building the Client

1. Navigate to the client directory:
   ```
   cd distributed-auth-system/client
   ```

2. Download dependencies:
   ```
   go mod tidy
   ```

3. Build the client:
   ```
   go build -o auth-client
   ```

## Running the Client

The client supports several operations:

### Storing Data

```
./auth-client -server=localhost:50051 -op=store -key=user123 -value=password456
```

### Retrieving Data

```
./auth-client -server=localhost:50051 -op=get -key=user123
```

### Simulating Authentication

```
./auth-client -server=localhost:50051 -op=auth -key=user123 -value=password456
```

## Command-Line Options

- `-server`: Server address (default: "localhost:50051")
- `-tls`: Use TLS (default: false)
- `-op`: Operation: store, get, or auth (default: "store")
- `-key`: Key for store/get operations (default: "user123")
- `-value`: Value for store operation (default: "password456")

## Example Usage Scenario

1. Start the distributed auth system servers (in separate terminals):
   ```
   # Terminal 1
   cd distributed-auth-system
   go run . -node=node1 -port=:50051

   # Terminal 2
   cd distributed-auth-system
   go run . -node=node2 -port=:50052

   # Terminal 3
   cd distributed-auth-system
   go run . -node=node3 -port=:50053
   ```

2. Store a user's credentials:
   ```
   ./auth-client -op=store -key=john.doe -value=secure123
   ```

3. Authenticate the user:
   ```
   ./auth-client -op=auth -key=john.doe -value=secure123
   ```

4. Retrieve the user's data:
   ```
   ./auth-client -op=get -key=john.doe
   ```

## Notes

- This is a simplified client for demonstration purposes.
- In a real-world scenario, you would want to implement proper authentication and authorization mechanisms.
- The current implementation uses a simple key-value store for user credentials, which is not secure for production use.
