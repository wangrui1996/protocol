package protocol

import (
	"errors"
	"io"
	"net"
)

// Client is used to create connection.
type Client interface {
	Name() string
	Addr() string
	HasKey() bool
	SetKey(key string)
	Handshake(underlay  net.Conn) (io.ReadWriteCloser, error)
}

// ClientCreator is a function to create client.
type ClientCreator func(name string) (Client, error)

var (
	clientMap = make(map[string]ClientCreator)
)

// RegisterClient is used to register a client.
func RegisterClient(name string, c ClientCreator) {
	clientMap[name] = c
}

// ClientFromURL calls the registered creator to create client.
// dialer is the default upstream dialer so cannot be nil, we can use Default when calling this function.
func ClientFromInfo(protocol string, name string) (Client, error) {

	c, ok := clientMap[protocol]
	if ok {
		return c(name)
	}
	return nil, errors.New("unknown client scheme '" + protocol + "'")
}