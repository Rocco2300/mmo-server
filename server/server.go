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
	Conn *net.UDPConn

	Mutex            *sync.Mutex
	PlayerData       []core.Player
	PlayerConnection *sync.Map

	FreeId int

	Sim Simulation
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

	err = server.Sim.Init()
	if err != nil {
		return err
	}

	server.Mutex = new(sync.Mutex)

	return nil
}

func (server *Server) close() {
	server.Sim.Clean()
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

		if request.Type == "Connect" {
			server.handleConnection(request, remoteAddr)
		}

		if request.Type == "Move" {
			server.handleMove(request)
		}

		if request.Type == "Disconnect" {
			server.handleDisconnect(request)
		}

		if request.Type == "ChatMessage" {
			server.handleChatMessage(request)
		}
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

	if len(buf) > 65536 {
		panic("message too large")
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
		Body: core.Connect{
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

	x := float32(2.0 * (id % 5))
	z := float32(-2.0 * (id / 5))
	pos := rl.NewVector3(x, 0, z)
	if _, ok := server.PlayerConnection.Load(id); !ok {
		server.PlayerConnection.Store(server.FreeId, &addr)

		player := core.Player{
			Id:       id,
			Position: pos,
		}

		server.Mutex.Lock()
		server.PlayerData = append(server.PlayerData, player)
		server.Mutex.Unlock()
	}
	server.FreeId++
}

func (server *Server) handleDisconnect(request core.Message) {
	id := request.Body.(core.Disconnect).Id

	if _, ok := server.PlayerConnection.Load(id); !ok {
		fmt.Println("malformed disconnect request")
		return
	}

	server.PlayerConnection.Delete(id)
	var i int
	for i, _ = range server.PlayerData {
		if server.PlayerData[i].Id == id {
			break
		}
	}

	server.PlayerData[i] = server.PlayerData[len(server.PlayerData)-1]
	server.PlayerData = server.PlayerData[:len(server.PlayerData)-1]

	server.broadcast(request)
}

func (server *Server) handleMove(request core.Message) {
	id := request.Body.(core.Move).Id
	vel := request.Body.(core.Move).Offset

	var maxLength float32 = 10.0
	vel = rl.Vector3Scale(rl.Vector3Normalize(vel), maxLength)

	for i, player := range server.PlayerData {
		if player.Id == id {
			server.Mutex.Lock()
			server.PlayerData[i].Velocity = vel
			server.Mutex.Unlock()
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

func (server *Server) handleChatMessage(request core.Message) {
	response := core.Message{
		Body: core.ChatMessage{
			Id:      request.Body.(core.ChatMessage).Id,
			Content: request.Body.(core.ChatMessage).Content,
		},
	}

	server.broadcast(response)
}
