package distributed_auth_system

// Transaction represents a key-value transaction with authentication and signature
type Transaction struct {
	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Key       string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value     string `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	Signature []byte `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	PublicKey string `protobuf:"bytes,5,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	Timestamp int64  `protobuf:"varint,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

// Reset implements the protobuf.Message interface
func (x *Transaction) Reset() {
	*x = Transaction{}
}

// String implements the protobuf.Message interface
func (x *Transaction) String() string {
	return "Transaction"
}

// ProtoMessage implements the protobuf.Message interface
func (*Transaction) ProtoMessage() {}

// SyncRequest represents a request to synchronize state with another node
type SyncRequest struct {
	NodeId       string `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	MerkleRoot   string `protobuf:"bytes,2,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"`
	BloomFilter  []byte `protobuf:"bytes,3,opt,name=bloom_filter,json=bloomFilter,proto3" json:"bloom_filter,omitempty"`
	BatchSize    int32  `protobuf:"varint,4,opt,name=batch_size,json=batchSize,proto3" json:"batch_size,omitempty"`
	Continuation string `protobuf:"bytes,5,opt,name=continuation,proto3" json:"continuation,omitempty"`
}

// Reset implements the protobuf.Message interface
func (x *SyncRequest) Reset() {
	*x = SyncRequest{}
}

// String implements the protobuf.Message interface
func (x *SyncRequest) String() string {
	return "SyncRequest"
}

// ProtoMessage implements the protobuf.Message interface
func (*SyncRequest) ProtoMessage() {}

// SyncResponse represents a response to a sync request
type SyncResponse struct {
	Success             bool           `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	MerkleRoot          string         `protobuf:"bytes,2,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"`
	MissingTransactions []*Transaction `protobuf:"bytes,3,rep,name=missing_transactions,json=missingTransactions,proto3" json:"missing_transactions,omitempty"`
	ContinuationToken   string         `protobuf:"bytes,4,opt,name=continuation_token,json=continuationToken,proto3" json:"continuation_token,omitempty"`
	HasMore             bool           `protobuf:"varint,5,opt,name=has_more,json=hasMore,proto3" json:"has_more,omitempty"`
}

// Reset implements the protobuf.Message interface
func (x *SyncResponse) Reset() {
	*x = SyncResponse{}
}

// String implements the protobuf.Message interface
func (x *SyncResponse) String() string {
	return "SyncResponse"
}

// ProtoMessage implements the protobuf.Message interface
func (*SyncResponse) ProtoMessage() {}

// StateVerificationRequest represents a request to verify state with another node
type StateVerificationRequest struct {
	NodeId     string `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	MerkleRoot string `protobuf:"bytes,2,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"`
}

// Reset implements the protobuf.Message interface
func (x *StateVerificationRequest) Reset() {
	*x = StateVerificationRequest{}
}

// String implements the protobuf.Message interface
func (x *StateVerificationRequest) String() string {
	return "StateVerificationRequest"
}

// ProtoMessage implements the protobuf.Message interface
func (*StateVerificationRequest) ProtoMessage() {}

// StateVerificationResponse represents a response to a state verification request
type StateVerificationResponse struct {
	Consistent bool   `protobuf:"varint,1,opt,name=consistent,proto3" json:"consistent,omitempty"`
	MerkleRoot string `protobuf:"bytes,2,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"`
}

// Reset implements the protobuf.Message interface
func (x *StateVerificationResponse) Reset() {
	*x = StateVerificationResponse{}
}

// String implements the protobuf.Message interface
func (x *StateVerificationResponse) String() string {
	return "StateVerificationResponse"
}

// ProtoMessage implements the protobuf.Message interface
func (*StateVerificationResponse) ProtoMessage() {}
