package connection

import (
	"bufio"
	"net"
	"strconv"

	"github.com/recraft/recraft-lib/utils"
)

// Client definition
// This manages the low level connection between client and server
type Client struct {
	// Raw socket connection
	socket *net.TCPConn
	//Reader buffer
	Reader *bufio.Reader
	// Connected to the server?
	Connected bool
	// Address of the server
	Address string
	// Port of the server
	Port int16
}

func (client *Client) Read(data []byte) (int, error) {
	return client.socket.Read(data)
}

// Connect to a tcp server
func (client *Client) Connect() error {
	if client.Connected {
		return utils.NewError("Already connected")
	}
	address, err := net.ResolveTCPAddr("tcp", client.Address+":"+strconv.Itoa(int(client.Port)))
	if err != nil {
		return utils.NewError("Failed resoluting address")
	}
	client.socket, err = net.DialTCP("tcp", nil, address)
	if err != nil {
		return err
	}
	client.Connected = true
	client.Reader = bufio.NewReader(client.socket)
	return nil

}

// Send data
func (client *Client) Send(data []byte) (int, error) {
	return client.socket.Write(data)
}

// Close the current connection
func (client *Client) Close() error {
	err := client.socket.Close()
	if err != nil {
		return err
	}
	client.Reader = nil
	client.Connected = false
	return nil
}
