package main

import (
	"encoding/json"
	"fmt"
	"net"

	"mmo-server.local/core"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Client struct {
	Id   int
	Conn *net.UDPConn
}

func (client *Client) init() {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:12345")
	if err != nil {
		panic(err)
	}

	client.Conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
}

func (client *Client) close() {
	client.Conn.Close()
}

func (client *Client) write(message core.Message) {
	buf, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	_, err = client.Conn.Write(buf)
	if err != nil {
		panic(err)
	}
}

func (client *Client) read() core.Message {
	buf := make([]byte, 1024)
	n, _, err := client.Conn.ReadFromUDP(buf)
	if err != nil {
		panic(err)
	}

	var receivedMessage core.Message
	err = json.Unmarshal(buf[:n], &receivedMessage)
	if err != nil {
		panic(err)
	}

	return receivedMessage
}

func (client *Client) connect() {
	request := core.Message{
		Body: core.Connection{
			Id: -1,
		},
	}

	client.write(request)

	response := client.read()

	if response.Type == "Connection" {
		client.Id = response.Body.(core.Connection).Id
	}

	fmt.Printf("conected with id %d\n", client.Id)
}

func (client *Client) move(vel rl.Vector3) rl.Vector3 {
	validVel := rl.NewVector3(0, 0, 0)

	request := core.Message{
		Body: core.Move{
			Id:     client.Id,
			Offset: vel,
		},
	}

	client.write(request)

	response := client.read()

	if response.Body != nil {
		validVel = response.Body.(core.Move).Offset
	}

	return validVel
}
