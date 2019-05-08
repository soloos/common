package base

type DBOptionDBSWALTopicClusterItem struct {
	PeerID   string
	IsLeader int
}

type DBOptionDBSWALTopic struct {
	TopicID    string
	SWALCluter []DBOptionDBSWALTopicClusterItem
}

type DBOptionDBSDFS struct {
	Path                  string
	DefaultNetBlockCap    int
	DefaultMemBlockCap    int
	DefaultMemBlocksLimit int32
}

type DBOptions struct {
	DBSDFS      DBOptionDBSDFS
	DBSWALTopic DBOptionDBSWALTopic
}
