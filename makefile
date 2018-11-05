COMMON_LDFLAGS=""
COMMON_PREFIX=""

all:sdfsd

sdfsd:
	$(COMMON_PREFIX) go build -i -ldflags '$(COMMON_LDFLAGS)' -o ./bin/sdfsd sdfsd

include ./make/test
include ./make/bench

.PHONY:all soloos-server test
