package main

import (
	"encoding/json"
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
		Body: core.Connection{},
	}

	message, err := json.Marshal(messageStruct)
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(message)
	if err != nil {
		panic(err)
	}

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

		if rl.IsKeyDown(rl.KeyA) {
			pos.X -= 10.0 * dt
		}
		if rl.IsKeyDown(rl.KeyD) {
			pos.X += 10.0 * dt
		}
		if rl.IsKeyDown(rl.KeyS) {
			pos.Z += 10.0 * dt
		}
		if rl.IsKeyDown(rl.KeyW) {
			pos.Z -= 10.0 * dt
		}

		rl.BeginDrawing()

		rl.BeginMode3D(camera)

		rl.ClearBackground(rl.White)

		rl.DrawSphere(pos, 1, rl.Blue)
		rl.DrawGrid(10, 1)

		rl.EndMode3D()

		rl.EndDrawing()
	}
}
