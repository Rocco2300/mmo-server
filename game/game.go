package main

import (
	"encoding/json"
	"fmt"
	"net"

	rl "github.com/gen2brain/raylib-go/raylib"
	"mmo-server.local/core"
)

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:12345")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	messageStruct := core.Message{
		Body: core.Connection{
			Id: -1,
		},
	}

	message, err := json.Marshal(messageStruct)
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(message)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		panic(err)
	}

	var receivedMessage core.Message
	err = json.Unmarshal(buf[:n], &receivedMessage)
	if err != nil {
		panic(err)
	}

	var id int
	if receivedMessage.Type == "Connection" {
		id = receivedMessage.Body.(core.Connection).Id
	}

	fmt.Printf("conected with id %d\n", id)

	camera := rl.NewCamera3D(
		rl.NewVector3(0, 1, 3),
		rl.NewVector3(0, 0, 0),
		rl.NewVector3(0, 1, 0),
		45.0,
		rl.CameraPerspective,
	)

	pos := rl.NewVector3(0, 0, 0)

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		vel := rl.NewVector3(0, 0, 0)
		if rl.IsKeyDown(rl.KeyA) {
			vel.X -= 10.0 * dt
		}
		if rl.IsKeyDown(rl.KeyD) {
			vel.X += 10.0 * dt
		}
		if rl.IsKeyDown(rl.KeyS) {
			vel.Z += 10.0 * dt
		}
		if rl.IsKeyDown(rl.KeyW) {
			vel.Z -= 10.0 * dt
		}

		moveMessage := core.Message{
			Body: core.Move{
				Id:     id,
				Offset: vel,
			},
		}

		buf, err = json.Marshal(moveMessage)
		if err != nil {
			fmt.Printf("couldn't serialize move message: %v", err)
		}

		_, err = conn.Write(buf)
		if err != nil {
			fmt.Printf("couldn't send move message: %v", err)
		}

		buf = make([]byte, 1024)
		n, _, err = conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("couldn't read move message: %v", err)
		}

		if n != 0 {
			var validMove core.Message
			err = json.Unmarshal(buf[:n], &validMove)
			if err != nil {
				fmt.Errorf("couldn't deserialize move message")
			}

			validVel := rl.NewVector3(0, 0, 0)
			if validMove.Body != nil {
				validVel = validMove.Body.(core.Move).Offset
			} else {
				fmt.Println("rateu")
			}

			pos = rl.Vector3Add(pos, validVel)
		}

		rl.BeginDrawing()

		rl.BeginMode3D(camera)

		rl.ClearBackground(rl.White)

		rl.DrawSphere(pos, 1, rl.Blue)
		rl.DrawGrid(10, 1)

		rl.EndMode3D()

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}
}
