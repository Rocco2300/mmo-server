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

	Players []core.Player

	Dead bool
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

	client.Players = make([]core.Player, 0)
}

func (client *Client) close() {
	client.Conn.Close()
}

func (client *Client) listen() {
	go client.listenLoop()
}

func (client *Client) listenLoop() {
	for {
		if client.Dead {
			break
		}

		message, err := client.read()
		if err != nil {
			fmt.Println("error in reading message from server")
			continue
		}

		if message.Type == "GameState" {
			playerList := message.Body.(core.GameState).Players

			client.Players = playerList
		}
	}
}

func (client *Client) write(message core.Message) error {
	buf, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = client.Conn.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) read() (core.Message, error) {
	buf := make([]byte, 1024)
	n, _, err := client.Conn.ReadFromUDP(buf)
	if err != nil {
		return core.Message{}, err
	}

	var receivedMessage core.Message
	err = json.Unmarshal(buf[:n], &receivedMessage)
	if err != nil {
		return core.Message{}, err
	}

	return receivedMessage, nil
}

func (client *Client) connect() {
	request := core.Message{
		Body: core.Connect{
			Id: -1,
		},
	}

	err := client.write(request)
	if err != nil {
		fmt.Println("could not send connection request \n", err)
		return
	}

	response, err := client.read()
	if err != nil {
		fmt.Println("could not read connection response from server \n", err)
		return
	}

	if response.Type == "Connect" {
		client.Id = response.Body.(core.Connect).Id
	}

	fmt.Printf("conected with id %d\n", client.Id)
}

func (client *Client) disconnect() {
	client.Dead = true

	request := core.Message{
		Body: core.Disconnect{
			Id: client.Id,
		},
	}

	err := client.write(request)
	if err != nil {
		fmt.Println("couldn't cleanly disconnect")
	}
}

func (client *Client) move(vel rl.Vector3) {
	request := core.Message{
		Body: core.Move{
			Id:     client.Id,
			Offset: vel,
		},
	}

	client.write(request)
}
