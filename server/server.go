package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"mmo-server.local/core"
)

func main() {
	addr := &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 12345, Zone: ""}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	remoteConns := new(sync.Map)

	for {
		buf := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}

		var message core.Message
		err = json.Unmarshal(buf[:n], &message)
		if err != nil {
			//panic(err)
			fmt.Println(err)
		}

		if message.Type == "Connection" {
			fmt.Println("Connection request received")
		}

		if _, ok := remoteConns.Load(remoteAddr); !ok {
			remoteConns.Store(remoteAddr.String(), &remoteAddr)
		}

		go func() {
			remoteConns.Range(func(key, value interface{}) bool {
				if *value.(*net.Addr) == remoteAddr {
					return true
				}

				if _, err := conn.WriteTo(buf, *value.(*net.Addr)); err != nil {
					remoteConns.Delete(key)

					return true
				}

				return true
			})
		}()
	}
}
