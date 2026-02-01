package core

import (
	"encoding/json"
	"fmt"
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Id       int
	Position rl.Vector3
}

type MessageBody interface{}

type Connect struct {
	Id int
}

type Disconnect struct {
	Id int
}

type Spawn struct {
	Id  int
	Pos rl.Vector3
}

type GameState struct {
	Players []Player
}

type Move struct {
	Id     int        `json:"id"`
	Offset rl.Vector3 `json:"offset"`
}

var _ MessageBody = Move{}

var knowImplementations = []MessageBody{
	Connect{},
	Disconnect{},
	Spawn{},
	Move{},
	GameState{},
}

type Message struct {
	Type string      `json:"type"`
	Body MessageBody `json:"body"`
}

func (m *Message) UnmarshalJSON(bytes []byte) error {
	var data struct {
		Type string          `json:"type"`
		Body json.RawMessage `json:"body"`
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	for _, knownImplementation := range knowImplementations {
		knownType := reflect.TypeOf(knownImplementation)
		if knownType.Name() == data.Type {
			target := reflect.New(knownType)

			if err := json.Unmarshal(data.Body, target.Interface()); err != nil {
				return err
			}

			m.Type = data.Type
			m.Body = target.Elem().Interface().(MessageBody)
			return nil
		}
	}

	return fmt.Errorf("unknown message type: %s", data.Type)
}

func (m Message) MarshalJSON() (bytes []byte, err error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Body any    `json:"body"`
	}{
		Type: reflect.TypeOf(m.Body).Name(),
		Body: m.Body,
	})
}
