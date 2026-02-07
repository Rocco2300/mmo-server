package main

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"mmo-server.local/core"
)

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

		server.Mutex.Unlock()

		message := core.Message{
			Body: core.GameState{
				Players: server.PlayerData,
			},
		}

		server.broadcast(message)

		// TODO: might not time it very well
		// should change in the future with a manually timed clock
		time.Sleep(16 * time.Millisecond)
	}
}
