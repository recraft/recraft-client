package recraftclient

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"net"
	"strconv"
	"strings"

	"github.com/recraft/recraft-lib/types"
	jsontypes "github.com/recraft/recraft-lib/types/json"

	"github.com/recraft/recraft-lib/utils"
)

//RecraftClient defines a connection to a server, it should be defined only one time.
type RecraftClient struct {
	socket    net.Conn
	reader    *bufio.Reader
	connected bool
	Address   string
	joined    bool
}

//Connect to a Minecraft server
func (client *RecraftClient) Connect() error {
	var err error
	if client.connected {
		return errors.New("The server is already connected")
	}

	client.socket, err = net.Dial("tcp", client.Address)
	if err != nil {
		return err
	}
	client.connected = true
	//Open reader
	client.reader = bufio.NewReader(client.socket)
	return nil

}

//GetServerInfo : Get info about the server
func (client *RecraftClient) GetServerInfo() (jsontypes.ServerInfo, error) {
	serverInfo := jsontypes.ServerInfo{}
	if !client.connected {
		return serverInfo, errors.New("Not connected to any server")
	}
	if client.joined {
		return serverInfo, errors.New("Client joined! Can't do this request")

	}

	//Handshake structure
	handshake := types.ClientHandshake{}

	//Protocol version of the client
	handshake.ProtocolVersion = types.VarInt(utils.ProtocolVersion)
	//Address where the client originally connected
	handshake.HostAddress = types.String(strings.Split(client.Address, ":")[0])
	//Port where the client originally connected
	num, err := strconv.Atoi(strings.Split(client.Address, ":")[1])
	if err != nil {
		client.socket.Close()
		client.connected = false
		return serverInfo, err
	}
	handshake.Port = types.Short(num)

	//Next State (Status)
	handshake.NextState = utils.StatusState

	//Write to host
	Mcc, err := utils.StructToBinary(&handshake, 0)
	if err != nil {
		client.socket.Close()
		client.connected = false
		return serverInfo, err
	}
	_, err = client.socket.Write(Mcc)
	if err != nil {
		client.socket.Close()
		client.connected = false
		return serverInfo, err
	}
	client.socket.Write([]byte{1, 0})

	//Read the response's lenght
	var lenght types.VarInt
	err = lenght.Read(client.reader)
	if err != nil {
		client.socket.Close()
		client.connected = false
		return serverInfo, err
	}

	//Create a "sandboxed" response
	sandboxResponse, err := types.ReadBytes(client.reader, int(lenght))
	if err != nil {
		client.socket.Close()
		client.connected = false
		return serverInfo, err
	}
	sandboxReader := bufio.NewReader(bytes.NewReader(sandboxResponse))

	var packetID types.VarInt
	//Read packetid's value
	err = packetID.Read(sandboxReader)
	if err != nil {
		client.socket.Close()
		client.connected = false
		return serverInfo, err
	}

	if packetID == 0 {

		//Get server's response
		info := &types.ServerListPingResponse{}
		err = utils.BinaryToStruct(info, sandboxReader)
		if err != nil {
			client.socket.Close()
			client.connected = false
			return serverInfo, err
		}
		ping := types.Pong{Payload: types.Long(rand.Int63())}
		payload, err := utils.StructToBinary(&ping, 1)
		if err != nil {
			client.socket.Close()
			client.connected = false
			return serverInfo, err
		}
		client.socket.Write(payload)

		client.socket.Close()
		client.connected = false
		json.Unmarshal([]byte(info.JSON), &serverInfo)
		return serverInfo, nil

	}
	return serverInfo, errors.New("invalid packetid")

	/*
		//"Ping" command, Needed to receive the response
		client.socket.Write([]byte{1, 0})
		ruffer, err := utils.GetVarIntFromBuffer(client.reader)
		//Get the Packetid
		packetID, _ := utils.GetVarIntFromBuffer(client.reader)
		Buffer := make([]byte, 0)

		for i := int32(0); i < (int32(ruffer) - 1); i++ {
			read, err := client.reader.ReadByte()
			if err != nil {
				fmt.Println("error! ", err)
				return jsontypes.ServerInfo{}, err
			}
			Buffer = append(Buffer, read)
		}

		//fmt.Println("reading buffer")

		//Read the response
		//_, err = client.reader.Read(Buffer)

		if err != nil {
			client.socket.Close()
			client.connected = false
			return serverInfo, err
		}
		//Create a "BufferedIndexed" element
		CountBuffer := &utils.BufferIndexed{CurrentIndex: 0, Buffer: Buffer}

		//CountBuffer.GetVarInt()
		var errs error = nil
		if packetID == 0 {
			response, err := CountBuffer.ReadStringAsBytes()
			if err != nil {
				client.socket.Close()
				client.connected = false
				return serverInfo, err
			}
			json.Unmarshal(response, &serverInfo)

		} else {
			errs = errors.New("Invalid packetid: ")
		}*/
	//Another "Ping" request, this time the connection will be closed.
	/*_, err = client.socket.Write([]byte{1, 0})
	if err != nil {
		client.socket.Close()
		client.connected = false
		return serverInfo, err
	}
	client.connected = false
	client.socket.Close()
	return serverInfo, errs*/
}
