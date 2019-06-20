COMMON_LDFLAGS=""
COMMON_PREFIX=""

GOBUILD = $(SDFS_PREFIX) go build

clean-test-cache:
	go clean -testcache

fbs:
	flatc -o ./ -g ./snetprotocol/test.fbs
	flatc -o ./ -g ./sdfsprotocol/sdfs.fbs
	flatc -o ./ -g ./swalprotocol/swal.fbs

soloos-tool:
	$(GOBUILD) -o ./bin/soloos-tool soloos-tool

go-sql-parser:
	@cd ./3rdlib/github.com/pingcap/parser && \
		GO111MODULE=on go get -u github.com/pingcap/parser@master && \
	       	make parser

include ./make/test
include ./make/bench

.PHONY:all soloos-server test
