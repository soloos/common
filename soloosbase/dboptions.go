package soloosbase

type DBOptionDBSolomqTopicClusterItem struct {
	PeerID string
	Role   int
}

type DBOptionDBSolomq struct {
	TopicID            string
	SolomqCluter         []DBOptionDBSolomqTopicClusterItem
	DefaultNetBlockCap int
	DefaultMemBlockCap int
}

type DBOptionDBSolofs struct {
	NameSpaceID           int
	DefaultNetBlockCap    int
	DefaultMemBlockCap    int
	DefaultMemBlocksLimit int32
}

type DBOptions struct {
	DBSolofs DBOptionDBSolofs
	DBSolomq DBOptionDBSolomq
}
