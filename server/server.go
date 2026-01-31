package main

import (
	"net"
	"sync"
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
		_, remoteAddr, err := conn.ReadFrom(buf)
		if err != nil {
			continue
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
