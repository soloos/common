package soloosbase

type DBOptionDBSOLOMQTopicClusterItem struct {
	PeerID string
	Role   int
}

type DBOptionDBSOLOMQ struct {
	TopicID            string
	SOLOMQCluter         []DBOptionDBSOLOMQTopicClusterItem
	DefaultNetBlockCap int
	DefaultMemBlockCap int
}

type DBOptionDBSOLOFS struct {
	NameSpaceID           int
	DefaultNetBlockCap    int
	DefaultMemBlockCap    int
	DefaultMemBlocksLimit int32
}

type DBOptions struct {
	DBSOLOFS DBOptionDBSOLOFS
	DBSOLOMQ DBOptionDBSOLOMQ
}
