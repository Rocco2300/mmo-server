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
		message, err := client.read()
		if err != nil {
			fmt.Println("error in reading message from server")
			continue
		}

		if message.Type == "Spawn" {
			player := core.Player{
				Id:       message.Body.(core.Spawn).Id,
				Position: message.Body.(core.Spawn).Pos,
			}

			client.Players = append(client.Players, player)
			continue
		}

		if message.Type == "GameState" {
			playerList := message.Body.(core.GameState).Players

			for _, player := range playerList {
				client.Players = append(client.Players, player)
			}
		}

		if message.Type == "Move" {
			for i, player := range client.Players {
				if player.Id == message.Body.(core.Move).Id {
					velocity := message.Body.(core.Move).Offset
					client.Players[i].Position = rl.Vector3Add(player.Position, velocity)
				}
			}
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
		Body: core.Connection{
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

	if response.Type == "Connection" {
		client.Id = response.Body.(core.Connection).Id
	}

	fmt.Printf("conected with id %d\n", client.Id)
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
