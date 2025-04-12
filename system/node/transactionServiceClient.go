package node

import (
	"context"

	"github.com/GooseFuse/distributed-auth-system/protoc"
	"google.golang.org/grpc"
)

type TransactionServiceClient struct {
	cc *grpc.ClientConn
}

func NewTransactionServiceClient(cc *grpc.ClientConn) *TransactionServiceClient {
	return &TransactionServiceClient{cc: cc}
}

func (c *TransactionServiceClient) VerifyState(ctx context.Context, in *protoc.StateVerificationRequest, opts ...grpc.CallOption) (svr *protoc.StateVerificationResponse, e error) {
	out := new(protoc.StateVerificationResponse)
	err := c.cc.Invoke(ctx, "/TransactionService/VerifyState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *TransactionServiceClient) HandleTransaction(ctx context.Context, in *protoc.TransactionRequest, opts ...grpc.CallOption) (tr *protoc.TransactionResponse, e error) {
	out := new(protoc.TransactionResponse)
	err := c.cc.Invoke(ctx, "/TransactionService/HandleTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
