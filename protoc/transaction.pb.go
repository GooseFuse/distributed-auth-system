// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: transaction.proto

package protoc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Transaction represents a key-value transaction with authentication and signature
type Transaction struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Key           string                 `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value         string                 `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	Signature     []byte                 `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	PublicKey     string                 `protobuf:"bytes,5,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	Timestamp     int64                  `protobuf:"varint,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	mi := &file_transaction_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{0}
}

func (x *Transaction) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Transaction) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Transaction) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Transaction) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *Transaction) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

func (x *Transaction) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

// TransactionRequest represents a request to handle a transaction
type TransactionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Transaction   *Transaction           `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransactionRequest) Reset() {
	*x = TransactionRequest{}
	mi := &file_transaction_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransactionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionRequest) ProtoMessage() {}

func (x *TransactionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionRequest.ProtoReflect.Descriptor instead.
func (*TransactionRequest) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{1}
}

func (x *TransactionRequest) GetTransaction() *Transaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

// TransactionResponse represents a response from handling a transaction
type TransactionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage  string                 `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	MerkleRoot    string                 `protobuf:"bytes,3,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"` // Current Merkle root after transaction
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransactionResponse) Reset() {
	*x = TransactionResponse{}
	mi := &file_transaction_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionResponse) ProtoMessage() {}

func (x *TransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionResponse.ProtoReflect.Descriptor instead.
func (*TransactionResponse) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{2}
}

func (x *TransactionResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *TransactionResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

func (x *TransactionResponse) GetMerkleRoot() string {
	if x != nil {
		return x.MerkleRoot
	}
	return ""
}

// AppendEntriesRequest represents a request for the AppendEntries RPC in Raft
type AppendEntriesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Term          int32                  `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId      string                 `protobuf:"bytes,2,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	PrevLogIndex  int32                  `protobuf:"varint,3,opt,name=prev_log_index,json=prevLogIndex,proto3" json:"prev_log_index,omitempty"`
	PrevLogTerm   int32                  `protobuf:"varint,4,opt,name=prev_log_term,json=prevLogTerm,proto3" json:"prev_log_term,omitempty"`
	Entries       []*Transaction         `protobuf:"bytes,5,rep,name=entries,proto3" json:"entries,omitempty"`
	LeaderCommit  int32                  `protobuf:"varint,6,opt,name=leader_commit,json=leaderCommit,proto3" json:"leader_commit,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AppendEntriesRequest) Reset() {
	*x = AppendEntriesRequest{}
	mi := &file_transaction_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AppendEntriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesRequest) ProtoMessage() {}

func (x *AppendEntriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesRequest.ProtoReflect.Descriptor instead.
func (*AppendEntriesRequest) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{3}
}

func (x *AppendEntriesRequest) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntriesRequest) GetLeaderId() string {
	if x != nil {
		return x.LeaderId
	}
	return ""
}

func (x *AppendEntriesRequest) GetPrevLogIndex() int32 {
	if x != nil {
		return x.PrevLogIndex
	}
	return 0
}

func (x *AppendEntriesRequest) GetPrevLogTerm() int32 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

func (x *AppendEntriesRequest) GetEntries() []*Transaction {
	if x != nil {
		return x.Entries
	}
	return nil
}

func (x *AppendEntriesRequest) GetLeaderCommit() int32 {
	if x != nil {
		return x.LeaderCommit
	}
	return 0
}

// AppendEntriesResponse represents a response to the AppendEntries RPC
type AppendEntriesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Term          int32                  `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Success       bool                   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AppendEntriesResponse) Reset() {
	*x = AppendEntriesResponse{}
	mi := &file_transaction_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AppendEntriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesResponse) ProtoMessage() {}

func (x *AppendEntriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesResponse.ProtoReflect.Descriptor instead.
func (*AppendEntriesResponse) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{4}
}

func (x *AppendEntriesResponse) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntriesResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

// RequestVoteRequest represents a request for the RequestVote RPC in Raft
type RequestVoteRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Term          int32                  `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	CandidateId   string                 `protobuf:"bytes,2,opt,name=candidate_id,json=candidateId,proto3" json:"candidate_id,omitempty"`
	LastLogIndex  int32                  `protobuf:"varint,3,opt,name=last_log_index,json=lastLogIndex,proto3" json:"last_log_index,omitempty"`
	LastLogTerm   int32                  `protobuf:"varint,4,opt,name=last_log_term,json=lastLogTerm,proto3" json:"last_log_term,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RequestVoteRequest) Reset() {
	*x = RequestVoteRequest{}
	mi := &file_transaction_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RequestVoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVoteRequest) ProtoMessage() {}

func (x *RequestVoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVoteRequest.ProtoReflect.Descriptor instead.
func (*RequestVoteRequest) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{5}
}

func (x *RequestVoteRequest) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RequestVoteRequest) GetCandidateId() string {
	if x != nil {
		return x.CandidateId
	}
	return ""
}

func (x *RequestVoteRequest) GetLastLogIndex() int32 {
	if x != nil {
		return x.LastLogIndex
	}
	return 0
}

func (x *RequestVoteRequest) GetLastLogTerm() int32 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

// RequestVoteResponse represents a response to the RequestVote RPC
type RequestVoteResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Term          int32                  `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	VoteGranted   bool                   `protobuf:"varint,2,opt,name=vote_granted,json=voteGranted,proto3" json:"vote_granted,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RequestVoteResponse) Reset() {
	*x = RequestVoteResponse{}
	mi := &file_transaction_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RequestVoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVoteResponse) ProtoMessage() {}

func (x *RequestVoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVoteResponse.ProtoReflect.Descriptor instead.
func (*RequestVoteResponse) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{6}
}

func (x *RequestVoteResponse) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RequestVoteResponse) GetVoteGranted() bool {
	if x != nil {
		return x.VoteGranted
	}
	return false
}

// SyncRequest represents a request to synchronize state with another node
type SyncRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	NodeId        string                 `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	MerkleRoot    string                 `protobuf:"bytes,2,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"`
	BloomFilter   []byte                 `protobuf:"bytes,3,opt,name=bloom_filter,json=bloomFilter,proto3" json:"bloom_filter,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncRequest) Reset() {
	*x = SyncRequest{}
	mi := &file_transaction_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncRequest) ProtoMessage() {}

func (x *SyncRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncRequest.ProtoReflect.Descriptor instead.
func (*SyncRequest) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{7}
}

func (x *SyncRequest) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *SyncRequest) GetMerkleRoot() string {
	if x != nil {
		return x.MerkleRoot
	}
	return ""
}

func (x *SyncRequest) GetBloomFilter() []byte {
	if x != nil {
		return x.BloomFilter
	}
	return nil
}

// SyncResponse represents a response to a sync request
type SyncResponse struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	Success             bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	MerkleRoot          string                 `protobuf:"bytes,2,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"`
	MissingTransactions []*Transaction         `protobuf:"bytes,3,rep,name=missing_transactions,json=missingTransactions,proto3" json:"missing_transactions,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *SyncResponse) Reset() {
	*x = SyncResponse{}
	mi := &file_transaction_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncResponse) ProtoMessage() {}

func (x *SyncResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncResponse.ProtoReflect.Descriptor instead.
func (*SyncResponse) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{8}
}

func (x *SyncResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SyncResponse) GetMerkleRoot() string {
	if x != nil {
		return x.MerkleRoot
	}
	return ""
}

func (x *SyncResponse) GetMissingTransactions() []*Transaction {
	if x != nil {
		return x.MissingTransactions
	}
	return nil
}

// StateVerificationRequest represents a request to verify state with another node
type StateVerificationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	NodeId        string                 `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	MerkleRoot    string                 `protobuf:"bytes,2,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StateVerificationRequest) Reset() {
	*x = StateVerificationRequest{}
	mi := &file_transaction_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StateVerificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateVerificationRequest) ProtoMessage() {}

func (x *StateVerificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateVerificationRequest.ProtoReflect.Descriptor instead.
func (*StateVerificationRequest) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{9}
}

func (x *StateVerificationRequest) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *StateVerificationRequest) GetMerkleRoot() string {
	if x != nil {
		return x.MerkleRoot
	}
	return ""
}

// StateVerificationResponse represents a response to a state verification request
type StateVerificationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Consistent    bool                   `protobuf:"varint,1,opt,name=consistent,proto3" json:"consistent,omitempty"`
	MerkleRoot    string                 `protobuf:"bytes,2,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StateVerificationResponse) Reset() {
	*x = StateVerificationResponse{}
	mi := &file_transaction_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StateVerificationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateVerificationResponse) ProtoMessage() {}

func (x *StateVerificationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateVerificationResponse.ProtoReflect.Descriptor instead.
func (*StateVerificationResponse) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{10}
}

func (x *StateVerificationResponse) GetConsistent() bool {
	if x != nil {
		return x.Consistent
	}
	return false
}

func (x *StateVerificationResponse) GetMerkleRoot() string {
	if x != nil {
		return x.MerkleRoot
	}
	return ""
}

// NodeInfo represents information about a node in the network
type NodeInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Address       string                 `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	IsLeader      bool                   `protobuf:"varint,3,opt,name=is_leader,json=isLeader,proto3" json:"is_leader,omitempty"`
	Term          int32                  `protobuf:"varint,4,opt,name=term,proto3" json:"term,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NodeInfo) Reset() {
	*x = NodeInfo{}
	mi := &file_transaction_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NodeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeInfo) ProtoMessage() {}

func (x *NodeInfo) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeInfo.ProtoReflect.Descriptor instead.
func (*NodeInfo) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{11}
}

func (x *NodeInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *NodeInfo) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *NodeInfo) GetIsLeader() bool {
	if x != nil {
		return x.IsLeader
	}
	return false
}

func (x *NodeInfo) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

// NetworkInfoRequest represents a request for information about the network
type NetworkInfoRequest struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	RequestingNodeId string                 `protobuf:"bytes,1,opt,name=requesting_node_id,json=requestingNodeId,proto3" json:"requesting_node_id,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *NetworkInfoRequest) Reset() {
	*x = NetworkInfoRequest{}
	mi := &file_transaction_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NetworkInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkInfoRequest) ProtoMessage() {}

func (x *NetworkInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkInfoRequest.ProtoReflect.Descriptor instead.
func (*NetworkInfoRequest) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{12}
}

func (x *NetworkInfoRequest) GetRequestingNodeId() string {
	if x != nil {
		return x.RequestingNodeId
	}
	return ""
}

// NetworkInfoResponse represents a response with information about the network
type NetworkInfoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Nodes         []*NodeInfo            `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	LeaderId      string                 `protobuf:"bytes,2,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	CurrentTerm   int32                  `protobuf:"varint,3,opt,name=current_term,json=currentTerm,proto3" json:"current_term,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NetworkInfoResponse) Reset() {
	*x = NetworkInfoResponse{}
	mi := &file_transaction_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NetworkInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkInfoResponse) ProtoMessage() {}

func (x *NetworkInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkInfoResponse.ProtoReflect.Descriptor instead.
func (*NetworkInfoResponse) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{13}
}

func (x *NetworkInfoResponse) GetNodes() []*NodeInfo {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *NetworkInfoResponse) GetLeaderId() string {
	if x != nil {
		return x.LeaderId
	}
	return ""
}

func (x *NetworkInfoResponse) GetCurrentTerm() int32 {
	if x != nil {
		return x.CurrentTerm
	}
	return 0
}

// CheckpointRequest represents a request to create or restore from a checkpoint
type CheckpointRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	NodeId        string                 `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Create        bool                   `protobuf:"varint,2,opt,name=create,proto3" json:"create,omitempty"`                                 // true for create, false for restore
	CheckpointId  int64                  `protobuf:"varint,3,opt,name=checkpoint_id,json=checkpointId,proto3" json:"checkpoint_id,omitempty"` // Used for restore only
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckpointRequest) Reset() {
	*x = CheckpointRequest{}
	mi := &file_transaction_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckpointRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckpointRequest) ProtoMessage() {}

func (x *CheckpointRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckpointRequest.ProtoReflect.Descriptor instead.
func (*CheckpointRequest) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{14}
}

func (x *CheckpointRequest) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *CheckpointRequest) GetCreate() bool {
	if x != nil {
		return x.Create
	}
	return false
}

func (x *CheckpointRequest) GetCheckpointId() int64 {
	if x != nil {
		return x.CheckpointId
	}
	return 0
}

// CheckpointResponse represents a response to a checkpoint request
type CheckpointResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	CheckpointId  int64                  `protobuf:"varint,2,opt,name=checkpoint_id,json=checkpointId,proto3" json:"checkpoint_id,omitempty"`
	ErrorMessage  string                 `protobuf:"bytes,3,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckpointResponse) Reset() {
	*x = CheckpointResponse{}
	mi := &file_transaction_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckpointResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckpointResponse) ProtoMessage() {}

func (x *CheckpointResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckpointResponse.ProtoReflect.Descriptor instead.
func (*CheckpointResponse) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{15}
}

func (x *CheckpointResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *CheckpointResponse) GetCheckpointId() int64 {
	if x != nil {
		return x.CheckpointId
	}
	return 0
}

func (x *CheckpointResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

var File_transaction_proto protoreflect.FileDescriptor

const file_transaction_proto_rawDesc = "" +
	"\n" +
	"\x11transaction.proto\"\xa0\x01\n" +
	"\vTransaction\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x10\n" +
	"\x03key\x18\x02 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x03 \x01(\tR\x05value\x12\x1c\n" +
	"\tsignature\x18\x04 \x01(\fR\tsignature\x12\x1d\n" +
	"\n" +
	"public_key\x18\x05 \x01(\tR\tpublicKey\x12\x1c\n" +
	"\ttimestamp\x18\x06 \x01(\x03R\ttimestamp\"D\n" +
	"\x12TransactionRequest\x12.\n" +
	"\vtransaction\x18\x01 \x01(\v2\f.TransactionR\vtransaction\"u\n" +
	"\x13TransactionResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12#\n" +
	"\rerror_message\x18\x02 \x01(\tR\ferrorMessage\x12\x1f\n" +
	"\vmerkle_root\x18\x03 \x01(\tR\n" +
	"merkleRoot\"\xde\x01\n" +
	"\x14AppendEntriesRequest\x12\x12\n" +
	"\x04term\x18\x01 \x01(\x05R\x04term\x12\x1b\n" +
	"\tleader_id\x18\x02 \x01(\tR\bleaderId\x12$\n" +
	"\x0eprev_log_index\x18\x03 \x01(\x05R\fprevLogIndex\x12\"\n" +
	"\rprev_log_term\x18\x04 \x01(\x05R\vprevLogTerm\x12&\n" +
	"\aentries\x18\x05 \x03(\v2\f.TransactionR\aentries\x12#\n" +
	"\rleader_commit\x18\x06 \x01(\x05R\fleaderCommit\"E\n" +
	"\x15AppendEntriesResponse\x12\x12\n" +
	"\x04term\x18\x01 \x01(\x05R\x04term\x12\x18\n" +
	"\asuccess\x18\x02 \x01(\bR\asuccess\"\x95\x01\n" +
	"\x12RequestVoteRequest\x12\x12\n" +
	"\x04term\x18\x01 \x01(\x05R\x04term\x12!\n" +
	"\fcandidate_id\x18\x02 \x01(\tR\vcandidateId\x12$\n" +
	"\x0elast_log_index\x18\x03 \x01(\x05R\flastLogIndex\x12\"\n" +
	"\rlast_log_term\x18\x04 \x01(\x05R\vlastLogTerm\"L\n" +
	"\x13RequestVoteResponse\x12\x12\n" +
	"\x04term\x18\x01 \x01(\x05R\x04term\x12!\n" +
	"\fvote_granted\x18\x02 \x01(\bR\vvoteGranted\"j\n" +
	"\vSyncRequest\x12\x17\n" +
	"\anode_id\x18\x01 \x01(\tR\x06nodeId\x12\x1f\n" +
	"\vmerkle_root\x18\x02 \x01(\tR\n" +
	"merkleRoot\x12!\n" +
	"\fbloom_filter\x18\x03 \x01(\fR\vbloomFilter\"\x8a\x01\n" +
	"\fSyncResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x1f\n" +
	"\vmerkle_root\x18\x02 \x01(\tR\n" +
	"merkleRoot\x12?\n" +
	"\x14missing_transactions\x18\x03 \x03(\v2\f.TransactionR\x13missingTransactions\"T\n" +
	"\x18StateVerificationRequest\x12\x17\n" +
	"\anode_id\x18\x01 \x01(\tR\x06nodeId\x12\x1f\n" +
	"\vmerkle_root\x18\x02 \x01(\tR\n" +
	"merkleRoot\"\\\n" +
	"\x19StateVerificationResponse\x12\x1e\n" +
	"\n" +
	"consistent\x18\x01 \x01(\bR\n" +
	"consistent\x12\x1f\n" +
	"\vmerkle_root\x18\x02 \x01(\tR\n" +
	"merkleRoot\"e\n" +
	"\bNodeInfo\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x18\n" +
	"\aaddress\x18\x02 \x01(\tR\aaddress\x12\x1b\n" +
	"\tis_leader\x18\x03 \x01(\bR\bisLeader\x12\x12\n" +
	"\x04term\x18\x04 \x01(\x05R\x04term\"B\n" +
	"\x12NetworkInfoRequest\x12,\n" +
	"\x12requesting_node_id\x18\x01 \x01(\tR\x10requestingNodeId\"v\n" +
	"\x13NetworkInfoResponse\x12\x1f\n" +
	"\x05nodes\x18\x01 \x03(\v2\t.NodeInfoR\x05nodes\x12\x1b\n" +
	"\tleader_id\x18\x02 \x01(\tR\bleaderId\x12!\n" +
	"\fcurrent_term\x18\x03 \x01(\x05R\vcurrentTerm\"i\n" +
	"\x11CheckpointRequest\x12\x17\n" +
	"\anode_id\x18\x01 \x01(\tR\x06nodeId\x12\x16\n" +
	"\x06create\x18\x02 \x01(\bR\x06create\x12#\n" +
	"\rcheckpoint_id\x18\x03 \x01(\x03R\fcheckpointId\"x\n" +
	"\x12CheckpointResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12#\n" +
	"\rcheckpoint_id\x18\x02 \x01(\x03R\fcheckpointId\x12#\n" +
	"\rerror_message\x18\x03 \x01(\tR\ferrorMessage2\xb8\x03\n" +
	"\x12TransactionService\x12>\n" +
	"\x11HandleTransaction\x12\x13.TransactionRequest\x1a\x14.TransactionResponse\x12>\n" +
	"\rAppendEntries\x12\x15.AppendEntriesRequest\x1a\x16.AppendEntriesResponse\x128\n" +
	"\vRequestVote\x12\x13.RequestVoteRequest\x1a\x14.RequestVoteResponse\x12(\n" +
	"\tSyncState\x12\f.SyncRequest\x1a\r.SyncResponse\x12D\n" +
	"\vVerifyState\x12\x19.StateVerificationRequest\x1a\x1a.StateVerificationResponse\x12;\n" +
	"\x0eGetNetworkInfo\x12\x13.NetworkInfoRequest\x1a\x14.NetworkInfoResponse\x12;\n" +
	"\x10ManageCheckpoint\x12\x12.CheckpointRequest\x1a\x13.CheckpointResponseB5Z3github.com/GooseFuse/distributed-auth-system/protocb\x06proto3"

var (
	file_transaction_proto_rawDescOnce sync.Once
	file_transaction_proto_rawDescData []byte
)

func file_transaction_proto_rawDescGZIP() []byte {
	file_transaction_proto_rawDescOnce.Do(func() {
		file_transaction_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_transaction_proto_rawDesc), len(file_transaction_proto_rawDesc)))
	})
	return file_transaction_proto_rawDescData
}

var file_transaction_proto_msgTypes = make([]protoimpl.MessageInfo, 16)
var file_transaction_proto_goTypes = []any{
	(*Transaction)(nil),               // 0: Transaction
	(*TransactionRequest)(nil),        // 1: TransactionRequest
	(*TransactionResponse)(nil),       // 2: TransactionResponse
	(*AppendEntriesRequest)(nil),      // 3: AppendEntriesRequest
	(*AppendEntriesResponse)(nil),     // 4: AppendEntriesResponse
	(*RequestVoteRequest)(nil),        // 5: RequestVoteRequest
	(*RequestVoteResponse)(nil),       // 6: RequestVoteResponse
	(*SyncRequest)(nil),               // 7: SyncRequest
	(*SyncResponse)(nil),              // 8: SyncResponse
	(*StateVerificationRequest)(nil),  // 9: StateVerificationRequest
	(*StateVerificationResponse)(nil), // 10: StateVerificationResponse
	(*NodeInfo)(nil),                  // 11: NodeInfo
	(*NetworkInfoRequest)(nil),        // 12: NetworkInfoRequest
	(*NetworkInfoResponse)(nil),       // 13: NetworkInfoResponse
	(*CheckpointRequest)(nil),         // 14: CheckpointRequest
	(*CheckpointResponse)(nil),        // 15: CheckpointResponse
}
var file_transaction_proto_depIdxs = []int32{
	0,  // 0: TransactionRequest.transaction:type_name -> Transaction
	0,  // 1: AppendEntriesRequest.entries:type_name -> Transaction
	0,  // 2: SyncResponse.missing_transactions:type_name -> Transaction
	11, // 3: NetworkInfoResponse.nodes:type_name -> NodeInfo
	1,  // 4: TransactionService.HandleTransaction:input_type -> TransactionRequest
	3,  // 5: TransactionService.AppendEntries:input_type -> AppendEntriesRequest
	5,  // 6: TransactionService.RequestVote:input_type -> RequestVoteRequest
	7,  // 7: TransactionService.SyncState:input_type -> SyncRequest
	9,  // 8: TransactionService.VerifyState:input_type -> StateVerificationRequest
	12, // 9: TransactionService.GetNetworkInfo:input_type -> NetworkInfoRequest
	14, // 10: TransactionService.ManageCheckpoint:input_type -> CheckpointRequest
	2,  // 11: TransactionService.HandleTransaction:output_type -> TransactionResponse
	4,  // 12: TransactionService.AppendEntries:output_type -> AppendEntriesResponse
	6,  // 13: TransactionService.RequestVote:output_type -> RequestVoteResponse
	8,  // 14: TransactionService.SyncState:output_type -> SyncResponse
	10, // 15: TransactionService.VerifyState:output_type -> StateVerificationResponse
	13, // 16: TransactionService.GetNetworkInfo:output_type -> NetworkInfoResponse
	15, // 17: TransactionService.ManageCheckpoint:output_type -> CheckpointResponse
	11, // [11:18] is the sub-list for method output_type
	4,  // [4:11] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_transaction_proto_init() }
func file_transaction_proto_init() {
	if File_transaction_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_transaction_proto_rawDesc), len(file_transaction_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   16,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_transaction_proto_goTypes,
		DependencyIndexes: file_transaction_proto_depIdxs,
		MessageInfos:      file_transaction_proto_msgTypes,
	}.Build()
	File_transaction_proto = out.File
	file_transaction_proto_goTypes = nil
	file_transaction_proto_depIdxs = nil
}
