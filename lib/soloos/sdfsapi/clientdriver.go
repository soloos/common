package sdfsapi

type ClientDriver interface {
	Init(nameNodeSRPCServerAddr string,
		defaultMemBlockChunkSize int, defaultMemBlockChunksLimit int32,
		dbDriver string, dsn string,
	) error
	InitClient(client Client,
		defaultNetBlockCap int,
		defaultMemBlockCap int,
	) error
}
