module soloos/common

go 1.12

require (
	github.com/fatih/color v1.7.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gocraft/dbr v0.0.0-20190503023340-d3d1e2876df1
	github.com/google/flatbuffers v1.11.0
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726
	github.com/siddontang/go-mysql v0.0.0-20190618002340-dbe0224ac097
	github.com/stretchr/testify v1.3.0
	github.com/tinylib/msgp v1.1.0
	github.com/vChrysanthemum/iron v0.0.0-20190729144136-a55d70a4a7d3
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
	soloos/solodb v0.0.0
)

replace (
	soloos/common v0.0.0 => /soloos/common
	soloos/soloboat v0.0.0 => /soloos/soloboat
	soloos/solodb v0.0.0 => /soloos/solodb
	soloos/solofs v0.0.0 => /soloos/solofs
	soloos/solomq v0.0.0 => /soloos/solomq
)
