#!/bin/bash

echo "=== Client Demo ==="

echo "1. Storing user data..."
./auth-client -server=auth-node-1:6333 -op=store -key=john.doe -value=secure123

echo
echo "2. Retrieving user data..."
./auth-client -server=auth-node-2:6333 -op=get -key=john.doe

echo
echo "3. Authenticating user..."
./auth-client -server=auth-node-3:6333 -op=auth -key=john.doe -value=secure123