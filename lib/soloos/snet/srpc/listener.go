package srpc

import (
	"net"
)

func makeListener(network, address string) (ln net.Listener, err error) {
	return net.Listen(network, address)
}
