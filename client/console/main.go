package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"mmo-server.local/client"
)

func main() {
	clients := make([]client.Client, 0)

	rl.SetConfigFlags(rl.FlagWindowHidden)
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	fmt.Println("enter a command (type help for more info): ")

	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("unexpected input try again")
			fmt.Println(err)
			continue
		}

		tokens := strings.Split(line, " ")
		command := strings.TrimSpace(tokens[0])
		switch command {
		default:
		case "help":
			fmt.Println("commands: ")
			fmt.Println("help - show help")
			fmt.Println("exit - exit program and disconnect all clients")
			fmt.Println("connect <number> - connect (number) instances")
			fmt.Println("move <id> <x> <y> <z> - move (id) client by (x, y, z) offset")
			break
		case "connect":
			if len(tokens) < 2 {
				fmt.Println("invalid argument, see help")
				break
			}

			number, err := strconv.Atoi(strings.TrimSpace(tokens[1]))
			if err != nil {
				fmt.Println("invalid argument, see help")
				break
			}

			for i := 0; i < number; i++ {
				var cli client.Client
				cli.Init()

				cli.Connect()
				cli.Listen()

				clients = append(clients, cli)
			}
			break

		case "move":
			if len(tokens) < 5 {
				fmt.Println("invalid argument, see help")
				break
			}

			id, err := strconv.Atoi(strings.TrimSpace(tokens[1]))
			if err != nil || id < 0 || id >= len(clients) {
				fmt.Println("invalid argument, see help")
				break
			}

			x, err1 := strconv.Atoi(strings.TrimSpace(tokens[2]))
			y, err2 := strconv.Atoi(strings.TrimSpace(tokens[3]))
			z, err3 := strconv.Atoi(strings.TrimSpace(tokens[4]))
			if err1 != nil || err2 != nil || err3 != nil {
				fmt.Println("invalid argument, see help")
				break
			}

			offset := rl.NewVector3(float32(x), float32(y), float32(z))
			clients[id].Move(offset)
			break

		case "exit":
			for i, _ := range clients {
				clients[i].Disconnect()
				clients[i].Close()
			}

			rl.CloseWindow()
			os.Exit(0)
		}
	}
}
