COMMON_LDFLAGS=""
COMMON_PREFIX=""

SOLOOS_SNET_PROTOS = $(shell find lib/soloos/snet -name '*.fbs')
GENERATED_PROTOS = $(shell find lib/soloos/snet -name "*.fbs"| sed 's/\.fbs/\.fbs\.go/g')
SOURCES = $(shell find . -name '*.go') $(GENERATED_PROTOS)

GOBUILD = $(SDFS_PREFIX) go build

%.fbs.go: $(SOLOOS_SNET_PROTOS)
	flatc -o ./lib/soloos/snet -g $(SOLOOS_SNET_PROTOS)

fbs: $(GENERATED_PROTOS)

go-sql-parser:
	@cd ./3rdlib/github.com/pingcap/parser && \
		GO111MODULE=on go get -u github.com/pingcap/parser@master && \
	       	make parser

include ./make/test
include ./make/bench

.PHONY:all soloos-server test
