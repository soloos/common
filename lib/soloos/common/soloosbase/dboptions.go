package soloosbase

type DBOptionDBSWALTopicClusterItem struct {
	PeerID string
	Role   int
}

type DBOptionDBSWAL struct {
	TopicID            string
	SWALCluter         []DBOptionDBSWALTopicClusterItem
	DefaultNetBlockCap int
	DefaultMemBlockCap int
}

type DBOptionDBSDFS struct {
	DefaultNetBlockCap    int
	DefaultMemBlockCap    int
	DefaultMemBlocksLimit int32
}

type DBOptions struct {
	DBSDFS DBOptionDBSDFS
	DBSWAL DBOptionDBSWAL
}
