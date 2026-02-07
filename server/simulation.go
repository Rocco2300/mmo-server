package main

import (
	"errors"
	"fmt"
	"math"
	"syscall"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Simulation struct {
	physics syscall.Handle

	init   uintptr
	update uintptr
	clean  uintptr

	setData uintptr
	getData uintptr

	Count      int
	Positions  []rl.Vector3
	Velocities []rl.Vector3
}

func (s *Simulation) Init() error {
	s.physics, _ = syscall.LoadLibrary("C:/Users/grigo/Repos/particle-simulation/build-release/libphysics.dll")
	if s.physics == 0 {
		return errors.New("libphysics library not found")
	}

	s.init, _ = syscall.GetProcAddress(s.physics, "init")
	if s.init == 0 {
		return errors.New("error finding init")
	}

	s.update, _ = syscall.GetProcAddress(s.physics, "update")
	if s.update == 0 {
		return errors.New("error finding update")
	}

	s.clean, _ = syscall.GetProcAddress(s.physics, "clean")
	if s.clean == 0 {
		return errors.New("error finding clean")
	}

	s.setData, _ = syscall.GetProcAddress(s.physics, "setData")
	if s.setData == 0 {
		return errors.New("error finding setData")
	}

	s.getData, _ = syscall.GetProcAddress(s.physics, "getData")
	if s.getData == 0 {
		return errors.New("error finding getData")
	}

	_, _, retCall := syscall.SyscallN(s.init)
	if retCall != 0 {
		return errors.New("error init")
	}

	return nil
}

func (s *Simulation) Update(deltaTime float32) {
	s.SetData()

	_, _, callErr := syscall.SyscallN(s.update, uintptr(math.Float32bits(deltaTime)))
	if callErr != 0 {
		fmt.Println("error calling update function")
	}

	s.GetData()
}

func (s *Simulation) Clean() {
	_, _, retCall := syscall.SyscallN(s.clean)
	if retCall != 0 {
		fmt.Println("error clean")
		return
	}
}

func (s *Simulation) SetData() {
	if s.Positions == nil || len(s.Positions) == 0 || s.Velocities == nil || len(s.Velocities) == 0 {
		return
	}

	count := uintptr(s.Count)
	positions := uintptr(unsafe.Pointer(&s.Positions[0]))
	velocities := uintptr(unsafe.Pointer(&s.Velocities[0]))
	_, _, callErr := syscall.SyscallN(s.setData, positions, velocities, count)
	if callErr != 0 {
		fmt.Println("error calling setData")
		return
	}
}

func (s *Simulation) GetData() {
	if s.Positions == nil || len(s.Positions) == 0 || s.Velocities == nil || len(s.Velocities) == 0 {
		return
	}

	positions := uintptr(unsafe.Pointer(&s.Positions[0]))
	_, _, callErr := syscall.SyscallN(s.getData, positions)
	if callErr != 0 {
		fmt.Println("error calling getData")
		return
	}
}
