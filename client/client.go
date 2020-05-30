package client

import "github.com/recraft/recraft-client/connection"

//Player datas
type Player struct {
	connection      *connection.Client
	protocolVersion int
	joined          bool
}

//NewClient creates Player structure
func NewClient(address string, port int16) *Player {

	return &Player{

		connection: &connection.Client{
			Address: address,
			Port:    port,
		},
	}

}
