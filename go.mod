module soloos/common

go 1.12

require (
	github.com/fatih/color v1.7.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gocraft/dbr v0.0.0-20190503023340-d3d1e2876df1
	github.com/google/flatbuffers v1.11.0
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726
	github.com/siddontang/go-mysql v0.0.0-20190618002340-dbe0224ac097
	github.com/stretchr/testify v1.3.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sys v0.0.0-20190606165138-5da285871e9c
	golang.org/x/xerrors v0.0.0-20190513163551-3ee3066db522
	google.golang.org/appengine v1.6.1 // indirect
	soloos/sdbone v0.0.0
)

replace (
	soloos/common v0.0.0 => /soloos/common
	soloos/sdbone v0.0.0 => /soloos/sdbone
)
