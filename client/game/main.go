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

	var prev int8 = 0
	for !rl.WindowShouldClose() {
		var curr int8 = 0
		offset := rl.NewVector3(0, 0, 0)
		if rl.IsKeyDown(rl.KeyA) {
			curr |= 0x01
			offset.X -= 10.0
		}
		if rl.IsKeyDown(rl.KeyD) {
			curr |= 0x02
			offset.X += 10.0
		}
		if rl.IsKeyDown(rl.KeyS) {
			curr |= 0x04
			offset.Z += 10.0
		}
		if rl.IsKeyDown(rl.KeyW) {
			curr |= 0x08
			offset.Z -= 10.0
		}

		if prev != curr {
			cli.Move(offset)
		}

		prev = curr

		rl.BeginDrawing()

		rl.BeginMode3D(camera)

		rl.ClearBackground(rl.White)

		// TODO: this might not be good due to race conditions
		i := 0
		for _, player := range cli.Players {
			if i > 10 {
				break
			}

			pos := player.Position

			rl.DrawSphere(pos, 1, rl.Blue)

			i++
		}

		rl.DrawGrid(100, 1)

		rl.EndMode3D()

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	cli.Disconnect()
}
