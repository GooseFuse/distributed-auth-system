package node

import "context"

type Client struct {
}

func (c *Client) VerifyState(ctx context.Context, req *StateVerificationRequest) (resp StateVerificationResponse, err error) {
	return resp, err
}
