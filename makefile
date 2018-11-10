COMMON_LDFLAGS=""
COMMON_PREFIX=""

SOLOOS_SNET_PROTOS = $(shell find protocol/soloos/snet -name '*.fbs')
GENERATED_PROTOS = $(shell find protocol/soloos -name "*.fbs"| sed 's/\.fbs/\.fbs\.go/g')
SOURCES = $(shell find . -name '*.go') $(GENERATED_PROTOS)

%.fbs.go: $(SOLOOS_SNET_PROTOS)
	flatc -o ./lib/soloos/snet -g $(SOLOOS_SNET_PROTOS)

fbs: $(GENERATED_PROTOS)

include ./make/test
include ./make/bench

.PHONY:all soloos-server test
