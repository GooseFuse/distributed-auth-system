package distributed_auth_system

import (
	"context"

	"google.golang.org/grpc"
)

// Constants for method names
const (
	TransactionService_SyncState_FullMethodName   = "/TransactionService/SyncState"
	TransactionService_VerifyState_FullMethodName = "/TransactionService/VerifyState"
)

// ExtendedTransactionServiceClient extends the generated TransactionServiceClient
// to include additional methods defined in the proto file but not generated
type ExtendedTransactionServiceClient interface {
	TransactionServiceClient
	SyncState(ctx context.Context, in *SyncRequest, opts ...grpc.CallOption) (*SyncResponse, error)
	VerifyState(ctx context.Context, in *StateVerificationRequest, opts ...grpc.CallOption) (*StateVerificationResponse, error)
}

// extendedTransactionServiceClient implements ExtendedTransactionServiceClient
type extendedTransactionServiceClient struct {
	TransactionServiceClient
	cc grpc.ClientConnInterface
}

// NewExtendedTransactionServiceClient creates a new ExtendedTransactionServiceClient
func NewExtendedTransactionServiceClient(cc grpc.ClientConnInterface) ExtendedTransactionServiceClient {
	return &extendedTransactionServiceClient{
		TransactionServiceClient: NewTransactionServiceClient(cc),
		cc:                       cc,
	}
}

// SyncState implements the SyncState RPC
func (c *extendedTransactionServiceClient) SyncState(ctx context.Context, in *SyncRequest, opts ...grpc.CallOption) (*SyncResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SyncResponse)
	err := c.cc.Invoke(ctx, TransactionService_SyncState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VerifyState implements the VerifyState RPC
func (c *extendedTransactionServiceClient) VerifyState(ctx context.Context, in *StateVerificationRequest, opts ...grpc.CallOption) (*StateVerificationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StateVerificationResponse)
	err := c.cc.Invoke(ctx, TransactionService_VerifyState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
