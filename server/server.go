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
	Conn    *net.UDPConn
	Players *sync.Map

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
	server.Players = new(sync.Map)

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
			server.handleConnection(remoteAddr, request)
		}

		if request.Type == "Move" {
			server.handleMove(request)
		}

		//go server.broadcast(request)
	}
}

func (server *Server) broadcast(message core.Message) {
	buf, err := json.Marshal(message)
	if err != nil {
		fmt.Errorf("error in converting message to json: %v", err)
		return
	}

	server.Players.Range(func(key, value interface{}) bool {
		if _, err := server.Conn.WriteTo(buf, *value.(*net.Addr)); err != nil {
			server.Players.Delete(key)

			return true
		}

		return true
	})
}

func (server *Server) handleConnection(addr net.Addr, request core.Message) {
	response := core.Message{
		Type: request.Type,
		Body: core.Connection{
			Id: server.FreeId,
		},
	}

	if _, ok := server.Players.Load(addr); !ok {
		server.Players.Store(server.FreeId, &addr)
	}
	server.FreeId++

	buffer, err := json.Marshal(response)
	if err != nil {
		fmt.Errorf("error in converting message to json: %v", err)
		return
	}

	_, err = server.Conn.WriteTo(buffer, addr)
	if err != nil {
		fmt.Errorf("error in sending connection message: %v", err)
	}

	fmt.Println("connection request received")
}

func (server *Server) handleMove(request core.Message) {
	id := request.Body.(core.Move).Id
	vel := request.Body.(core.Move).Offset

	var maxLength float32 = 10.0 * 0.016
	vel = rl.Vector3Scale(rl.Vector3Normalize(vel), maxLength)

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
