COMMON_LDFLAGS=""
COMMON_PREFIX=""

GOBUILD = $(SOLOFS_PREFIX) go build

clean-test-cache:
	go clean -testcache

protocol:
	go generate ./snetprotocol
	go generate ./solofsprotocol
	go generate ./solomqprotocol
	# flatc -o ./ -g ./snetprotocol/test.fbs
	# flatc -o ./ -g ./solofsprotocol/solofs.fbs
	# flatc -o ./ -g ./solomqprotocol/solomq.fbs

soloos-tool:
	$(GOBUILD) -o ./bin/soloos-tool soloos-tool

go-sql-parser:
	@cd ./3rdlib/github.com/pingcap/parser && \
		GO111MODULE=on go get -u github.com/pingcap/parser@master && \
	       	make parser

include ./make/test
include ./make/bench

.PHONY:all soloos-server test
