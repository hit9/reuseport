Reuseport
=========

Package reuseport provides SO_REUSEPORT in go.

Example
-------

TCP Listener:

```go
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
```

UDP PacketConn:

```go
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
```

Notes
-----

https://github.com/matishsiao/go_reuseport

https://github.com/kavu/go_reuseport

License
-------

MIT 2016 Chao Wang.
