package main

import (
	"bufio"
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	"mmo-server.local/client"
)

func messageLoop(client *client.Client) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("unexpected input try again")
			fmt.Println(err)
			continue
		}

		client.SendMessage(line)
	}
}

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var cli client.Client
	cli.Init()
	defer cli.Close()

	cli.Connect()

	cli.Listen()

	go messageLoop(&cli)

	camera := rl.NewCamera3D(
		rl.NewVector3(0, 10, 5),
		rl.NewVector3(0, 0, 0),
		rl.NewVector3(0, 1, 0),
		45.0,
		rl.CameraPerspective,
	)

	for !rl.WindowShouldClose() {
		offset := rl.NewVector3(0, 0, 0)
		if rl.IsKeyDown(rl.KeyA) {
			offset.X -= 10.0
		}
		if rl.IsKeyDown(rl.KeyD) {
			offset.X += 10.0
		}
		if rl.IsKeyDown(rl.KeyS) {
			offset.Z += 10.0
		}
		if rl.IsKeyDown(rl.KeyW) {
			offset.Z -= 10.0
		}

		cli.Move(offset)

		rl.BeginDrawing()

		rl.BeginMode3D(camera)

		rl.ClearBackground(rl.White)

		// TODO: this might not be good due to race conditions
		for _, player := range cli.Players {
			pos := player.Position

			rl.DrawSphere(pos, 1, rl.Blue)
		}

		rl.DrawGrid(100, 1)

		rl.EndMode3D()

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	cli.Disconnect()
}
