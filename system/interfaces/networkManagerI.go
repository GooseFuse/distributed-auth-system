package interfaces

import "github.com/GooseFuse/distributed-auth-system/protoc"

type NetworkManagerI interface {
	GetNodeId() string
	GetPeerClients() map[string]protoc.TransactionServiceClient
	GetPeerUrl(id string) string
}
