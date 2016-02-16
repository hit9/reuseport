// Copyright 2016 hit9. All rights reserved.

/*

Package reuseport provides SO_REUSEPORT in go.

TCP Listener

Example simple http server with reuseport enabled:

	package main

	import (
		"fmt"
		"github.com/hit9/reuseport"
		"net/http"
		"os"
	)

	func main() {
		ln, err := reuseport.Listener("tcp", ":8080")
		if err != nil {
			panic(err)
		}
		defer ln.Close()
		server := &http.Server{}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello from process %d\n", os.Getpid())
		})
		panic(server.Serve(ln))
	}

UDP Packet Conn

Example simple udp server with reuseport enabled:

	package main

	import (
		"fmt"
		"github.com/hit9/reuseport"
	)

	func main() {
		conn, err := reuseport.PacketConn("udp", ":8080")
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		buf := make([]byte, 2048)
		for {
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				break
			}
			fmt.Printf("Got message: %s\n", string(buf))
		}
	}

Ref And Thanks

https://github.com/matishsiao/go_reuseport

https://github.com/kavu/go_reuseport

*/
package reuseport
