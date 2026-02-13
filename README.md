# MMO Server

This project was developed for a university assignment. The requirement was a functioning server capable of handling movement and physics for multiple clients, along with a client-side interface. I used Go for its simplicity in serialization and networking. To reduce development time, I integrated a DLL-based physics API from the [Particle Simulation](https://github.com/Rocco2300/particle-simulation) project.

## Requirements

- go
- git
- CMake 3.30 or newer
- MinGW 15.2.0 or equivalent

## Building

```
git clone https://github.com/Rocco2300/mmo-server

cd mmo-server
./build.bat
```

## Usage

First you will have to start the server
```
./build/server.exe
```

Then you can connect as many players 
```
./build/game.exe
```

You can control the players using WASD.

Additionally, you can connect players en-masse using the console
```
./build/console.exe

connect 10
move 0 10 0 10

help

exit
```
