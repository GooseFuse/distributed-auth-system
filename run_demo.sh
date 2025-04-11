#!/bin/bash

# This script demonstrates how to run a cluster of distributed auth system nodes
# and interact with them using the client.

# Function to clean up background processes on exit
cleanup() {
    echo "Cleaning up..."
    kill $NODE1_PID $NODE2_PID $NODE3_PID 2>/dev/null
    exit
}

# Set up trap to catch Ctrl+C
trap cleanup SIGINT SIGTERM

# Build the system
echo "Building the distributed auth system..."
go build -o auth-system

# Build the client
echo "Building the client..."
cd client
go mod tidy
go build -o auth-client
cd ..

# Start Redis if not already running
echo "Checking if Redis is running..."
redis-cli ping > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "Starting Redis..."
    redis-server --daemonize yes
    sleep 2
fi

# Create nodes directory
mkdir -p data/nodes/node1 data/nodes/node2 data/nodes/node3

# Start three nodes in the background
echo "Starting node 1..."
./auth-system -node=node1 -port=:50051 -db=./data -redis=localhost:6379 &
NODE1_PID=$!
sleep 2

echo "Starting node 2..."
./auth-system -node=node2 -port=:50052 -db=./data -redis=localhost:6379 &
NODE2_PID=$!
sleep 2

echo "Starting node 3..."
./auth-system -node=node3 -port=:50053 -db=./data -redis=localhost:6379 &
NODE3_PID=$!
sleep 2

echo "All nodes are running!"
echo "Node 1 PID: $NODE1_PID"
echo "Node 2 PID: $NODE2_PID"
echo "Node 3 PID: $NODE3_PID"

# Demonstrate client operations
echo -e "\n=== Client Demo ==="
echo "1. Storing user data..."
./client/auth-client -server=localhost:50051 -op=store -key=john.doe -value=secure123

echo -e "\n2. Retrieving user data..."
./client/auth-client -server=localhost:50052 -op=get -key=john.doe

echo -e "\n3. Authenticating user..."
./client/auth-client -server=localhost:50053 -op=auth -key=john.doe -value=secure123

echo -e "\nDemo completed! The nodes are still running."
echo "Press Ctrl+C to stop all nodes and clean up."

# Wait for user to press Ctrl+C
wait
