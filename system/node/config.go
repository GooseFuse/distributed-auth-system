package node

type Config struct {
	NodeID             string
	Port               string
	DBPath             string
	RedisAddr          string
	PeerAddresses      map[string]string
	UseTLS             bool
	ElectionTimeoutMin int
	ElectionTimeoutMax int
	HeartbeatInterval  int
}
