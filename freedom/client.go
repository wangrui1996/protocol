package pure

import (
	"ethProxy/protocol"
	"io"
	"net"
)

func init() {
	protocol.RegisterClient(Name, NewFreedomClient)
}

func NewFreedomClient(name string) (protocol.Client, error) {
	c := &Client{name: name}
	return c, nil
}

func (c *Client) Name() string { return Name }

func (c *Client) Addr() string { return c.addr }

func (c *Client) HasKey() bool {
	return false
}

func (c *Client)SetKey(key string) {
	c.key = key
}

// Client is a vmess client
type Client struct {
	name string
	addr string
	key string
}

func (c *Client) Handshake(underlay net.Conn) (io.ReadWriteCloser, error) {
	return underlay, nil
}