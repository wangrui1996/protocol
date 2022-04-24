package sesp

import (
	"crypto/aes"
	"crypto/cipher"
	"ethProxy/protocol"
	"io"
	"math/rand"
	"net"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func init() {
	rand.Seed(time.Now().UnixNano())
	protocol.RegisterClient(Name, NewVmessClient)
}

func NewVmessClient(name string) (protocol.Client, error) {
	c := &Client{name: name}
	return c, nil
}
func (c *Client) Name() string { return Name }

func (c *Client) Addr() string { return c.addr }


func (c *Client) HasKey() bool {
	return true
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

	conn := &ClientConn{}
	conn.Conn = underlay
	copy(conn.BodyKey[:], c.key)
	randbytes := []byte(RandStringRunes(16))
	copy(conn.reqBodyIV[:], randbytes)
	// convert key
	//fmt.Printf("send: %s.\n", string(conn.reqBodyIV[:]))
	conn.Conn.Write(conn.reqBodyIV[:])
	//conn.Conn.Write([]byte("\n"))
	var respBodyIV [16]byte
	_, err := io.ReadFull(conn.Conn, respBodyIV[:])
	if err != nil {
		return nil, err
	}
	//fmt.Printf("get bytes: %s.", string(respBodyIV[:]))
	conn.respBodyIV = respBodyIV

	return conn, nil
}

// ClientConn is a connection to vmess server
type ClientConn struct {
	reqBodyIV   [16]byte
	//BodyKey  [16]byte
	respBodyIV  [16]byte
	BodyKey [16]byte

	net.Conn
	dataReader io.Reader
	dataWriter io.Writer
}



func (c *ClientConn) Write(b []byte) (n int, err error) {
	if c.dataWriter != nil {
		return c.dataWriter.Write(b)
	}

	c.dataWriter = c.Conn

	block, _ := aes.NewCipher(c.BodyKey[:])
	aead, _ := cipher.NewGCM(block)
	c.dataWriter = AEADWriter(c.Conn, aead, c.reqBodyIV[:])

	return c.dataWriter.Write(b)
}

func (c *ClientConn) Read(b []byte) (n int, err error) {
	if c.dataReader != nil {
		return c.dataReader.Read(b)
	}
	block, _ := aes.NewCipher(c.BodyKey[:])
	aead, _ := cipher.NewGCM(block)
	c.dataReader = AEADReader(c.Conn, aead, c.respBodyIV[:])
	return c.dataReader.Read(b)
}