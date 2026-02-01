package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var client Client
	client.init()
	defer client.close()

	client.connect()

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

		offset := rl.NewVector3(0, 0, 0)
		if rl.IsKeyDown(rl.KeyA) {
			offset.X -= 10.0 * dt
		}
		if rl.IsKeyDown(rl.KeyD) {
			offset.X += 10.0 * dt
		}
		if rl.IsKeyDown(rl.KeyS) {
			offset.Z += 10.0 * dt
		}
		if rl.IsKeyDown(rl.KeyW) {
			offset.Z -= 10.0 * dt
		}

		vel := client.move(offset)
		pos = rl.Vector3Add(pos, vel)

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
