package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"mmo-server.local/client"
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowHidden)
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	var cli client.Client
	cli.Init()
	defer cli.Close()

	cli.Connect()
	cli.Listen()
}
