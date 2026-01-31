package main

import (
	"encoding/json"
	"fmt"

	"mmo-server.local/core"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	m := core.Message{
		Body: core.Move{
			Id:     1,
			Offset: rl.NewVector3(1, 2, 3),
		},
	}

	data, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Cannot marshal message to json")
		panic(err)
	}

	var unmarshaled core.Message
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		fmt.Println("Cannot unmarshal message to json")
		panic(err)
	}

	switch v := unmarshaled.Body.(type) {
	case core.Move:
		fmt.Println("Message is move")
	default:
		panic(fmt.Sprintf("Unknown message type %T", v))
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
