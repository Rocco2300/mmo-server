package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"mmo-server.local/core"
)

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func main() {
	server := Server{}

	err := server.init()
	if err != nil {
		panic(err)
	}

	defer server.close()

	go server.listen()

	for {
		server.Mutex.Lock()
		server.Sim.Positions = make([]rl.Vector3, 0)
		server.Sim.Velocities = make([]rl.Vector3, 0)
		for _, player := range server.PlayerData {
			server.Sim.Positions = append(server.Sim.Positions, player.Position)
			server.Sim.Velocities = append(server.Sim.Velocities, player.Velocity)
		}
		server.Sim.Count = len(server.PlayerData)

		var deltaTime float32 = 0.0166
		server.Sim.Update(deltaTime)

		for i, _ := range server.PlayerData {
			server.PlayerData[i].Position = server.Sim.Positions[i]
		}

		i := 0
		for {
			isChunk := i > 0
			lowerBound := i * 100
			upperBound := min(len(server.PlayerData), (i+1)*100)

			if lowerBound >= len(server.PlayerData) {
				break
			}

			if i > 0 && len(server.PlayerData) == len(server.PlayerData[0:100]) {
				panic("what the hell")
			}

			fmt.Printf("sending from %d to %d\n", lowerBound, upperBound)

			players := append([]core.Player(nil), server.PlayerData[lowerBound:upperBound]...)
			fmt.Println("players out: ", len(players), cap(players))
			message := core.Message{
				Body: core.GameState{
					IsChunk: isChunk,
					Players: players,
				},
			}

			server.broadcast(message)
			i++
		}
		server.Mutex.Unlock()

		// TODO: might not time it very well
		// should change in the future with a manually timed clock
		time.Sleep(33 * time.Millisecond)
	}
}
