#!/bin/bash

sleep 30
echo "=== Client Demo ==="

echo "1. Storing user data..."
./auth-client -server=auth-node-1:6333 -op=store -key=john.doe -value=secure123

sleep 5

echo
echo "2. Retrieving user data..."
./auth-client -server=auth-node-2:6333 -op=get -key=john.doe

sleep 5

echo
echo "3. Authenticating user..."
./auth-client -server=auth-node-3:6333 -op=auth -key=john.doe -value=secure123


sleep 5

echo
echo "4. Authenticating user with bad creds..."
./auth-client -server=auth-node-3:6333 -op=auth -key=john.doe -value=secure1234