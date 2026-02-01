package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"mmo-server.local/core"
)

type Server struct {
	Conn             *net.UDPConn
	PlayerData       []core.Player
	PlayerConnection *sync.Map

	FreeId int
}

func (server *Server) init() error {
	addr := &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 12345, Zone: ""}

	var err error
	server.Conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	server.FreeId = 0
	server.PlayerData = make([]core.Player, 0)
	server.PlayerConnection = new(sync.Map)

	return nil
}

func (server *Server) close() {
	server.Conn.Close()
}

func (server *Server) listen() {
	for {
		buf := make([]byte, 1024)
		n, remoteAddr, err := server.Conn.ReadFrom(buf)
		if err != nil {
			continue
		}

		var request core.Message
		err = json.Unmarshal(buf[:n], &request)
		if err != nil {
			panic(err)
		}

		if request.Type == "Connection" {
			server.handleConnection(request, remoteAddr)
		}

		if request.Type == "Move" {
			server.handleMove(request)
		}

		//go server.broadcast(request)
	}
}

func (server *Server) write(message core.Message, addr net.Addr) error {
	buffer, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = server.Conn.WriteTo(buffer, addr)
	if err != nil {
		return err
	}

	return nil
}

func (server *Server) broadcast(message core.Message) {
	buf, err := json.Marshal(message)
	if err != nil {
		fmt.Errorf("error in converting message to json: %v", err)
		return
	}

	server.PlayerConnection.Range(func(key, value interface{}) bool {
		if _, err := server.Conn.WriteTo(buf, *value.(*net.Addr)); err != nil {
			server.PlayerConnection.Delete(key)

			return true
		}

		return true
	})
}

func (server *Server) handleConnection(request core.Message, addr net.Addr) {
	id := server.FreeId

	response := core.Message{
		Type: request.Type,
		Body: core.Connection{
			Id: id,
		},
	}

	// send connection accept message with id
	err := server.write(response, addr)
	if err != nil {
		fmt.Println("error writing connection response: ", err)
		return
	}

	fmt.Println("connection request received")

	response = core.Message{
		Body: core.GameState{
			Players: server.PlayerData,
		},
	}

	// send game state message, positions of already connected players
	err = server.write(response, addr)
	if err != nil {
		fmt.Println("error writing connection response: ", err)
		return
	}

	pos := rl.NewVector3(float32(2*id), 0, 0)
	if _, ok := server.PlayerConnection.Load(addr); !ok {
		server.PlayerConnection.Store(server.FreeId, &addr)

		player := core.Player{
			Id:       id,
			Position: pos,
		}
		server.PlayerData = append(server.PlayerData, player)
	}
	server.FreeId++

	// broadcast spawn event of new player to everyone
	response = core.Message{
		Body: core.Spawn{
			Id:  id,
			Pos: pos,
		},
	}
	server.broadcast(response)
}

func (server *Server) handleMove(request core.Message) {
	id := request.Body.(core.Move).Id
	vel := request.Body.(core.Move).Offset

	var maxLength float32 = 10.0 * 0.016
	vel = rl.Vector3Scale(rl.Vector3Normalize(vel), maxLength)

	for i, player := range server.PlayerData {
		if player.Id == id {
			position := player.Position
			newPosition := rl.Vector3Add(position, vel)
			server.PlayerData[i].Position = newPosition
		}
	}

	response := core.Message{
		Body: core.Move{
			Id:     id,
			Offset: vel,
		},
	}

	server.broadcast(response)
}

func main() {
	server := Server{}

	err := server.init()
	if err != nil {
		panic(err)
	}

	defer server.close()

	server.listen()
}
