module mmo-server.local/client/game

go 1.25.6


require (
	github.com/ebitengine/purego v0.9.1 // indirect
	github.com/gen2brain/raylib-go/raylib v0.55.1 // indirect
	golang.org/x/exp v0.0.0-20260112195511-716be5621a96 // indirect
	golang.org/x/sys v0.40.0 // indirect

    mmo-server.local/core v0.0.0
    mmo-server.local/client v0.0.0
)

replace mmo-server.local/client => ../
replace mmo-server.local/core => ../../core
