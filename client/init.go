package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/recraft/recraft-lib/packets"

	clientpackets "github.com/recraft/recraft-lib/packets/client"
	serverpackets "github.com/recraft/recraft-lib/packets/server"

	"github.com/recraft/recraft-lib/types"
	jsontypes "github.com/recraft/recraft-lib/types/json"
	"github.com/recraft/recraft-lib/utils"
)

func (player *Player) handshake(nextState types.VarInt) error {

	//Error if not connected
	if !player.connection.Connected {
		return utils.NewError("Client is not connected!")
	}

	// Handshake structure
	handshake := new(clientpackets.PacketHandshake)

	// Protocol version of the client
	handshake.ProtocolVersion = types.VarInt(player.protocolVersion)

	// Address where the client originally connected
	handshake.HostAddress = types.String(player.connection.Address)

	// Port where the client originally connected
	handshake.Port = types.Short(player.connection.Port)

	// Next State
	handshake.NextState = nextState

	// Convert struct to binary
	bin, err := utils.StructToBinary(handshake, 0)
	if err != nil {
		player.connection.Close()
		return err
	}

	// Send to stream
	_, err = player.connection.Send(bin)
	if err != nil {
		player.connection.Close()
		return err
	}

	return nil
}

// Status of the server (Server info)
// Note: This function will automatically close the connection after the request
func (player *Player) Status() (*jsontypes.ServerInfo, error) {

	if player.joined {
		return nil, utils.NewError("Client already joined")
	}

	err := player.connection.Connect()
	if err != nil {
		return nil, err
	}

	err = player.handshake(types.VarInt(packets.STATUS))
	if err != nil {
		return nil, err
	}

	statusRequest := &clientpackets.PacketStatusRequest{}
	// Make status request
	statusRq, err := utils.StructToBinary(nil, statusRequest.ID())
	if err != nil {
		return nil, err
	}
	_, err = player.connection.Send(statusRq)
	if err != nil {
		return nil, err
	}

	// Read response

	// Read lenght
	buffer := make([]byte, 1024)
	buflen, err := player.connection.Read(buffer)
	if err != nil {
		player.connection.Close()
		return nil, err
	}
	lbuffer := bufio.NewReader(bytes.NewReader(buffer))

	var lenght types.VarInt
	err = lenght.Read(lbuffer)
	if err != nil {
		player.connection.Close()
		return nil, err
	}
	fmt.Println("ok")
	var fullBytesBuffer []byte
	if buflen >= int(lenght) {
		fullBytesBuffer = buffer
	} else {
		fullBytesBuffer, err = types.ReadBytes(player.connection.Reader, int(lenght)-(buflen-2))
	}
	if err != nil {
		player.connection.Close()
		return nil, err
	}

	fullBuffer := bufio.NewReader(bytes.NewReader(append(buffer, fullBytesBuffer...)))

	response := &serverpackets.PacketStatus{}

	var packedID types.VarInt
	packedID.Read(fullBuffer)
	err = packedID.Read(fullBuffer)
	if err != nil {
		player.connection.Close()
		return nil, err
	}
	if packedID != types.VarInt(response.ID()) {
		player.connection.Close()
		return nil, utils.NewError("Wrong response ID")
	}
	err = utils.BinaryToStruct(response, fullBuffer)
	if err != nil {
		player.connection.Close()
		return nil, err
	}

	serverInfo := &jsontypes.ServerInfo{}
	json.Unmarshal([]byte(response.JSON), serverInfo)

	player.connection.Close()

	return serverInfo, nil
}
