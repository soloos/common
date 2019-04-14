package sdfsapi

type ClientDriver interface {
	Init(nameNodeSRPCServerAddr string,
		dbDriver string, dsn string,
	) error
	InitClient(client Client,
		defaultNetBlockCap int,
		defaultMemBlockCap int,
		defaultMemBlocksLimit int32,
	) error
}
