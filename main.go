// Should be able to receive and forward messages among clients
// > yo
// < whats up?
// > Nothing much!!
// < same here, hbu?
// > same ol' same

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"

	"time"
)

// list of connected clients
var clients = []net.Conn{}

func broadcast(msg string, from net.Conn) {
	for _, c := range clients {
		if c != from {
			c.Write([]byte(string(msg)))
		}
	}

}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	clients = append(clients, c)
	fmt.Println("appending client: ", len(clients))

	fmt.Printf("net.conn: %v\n", c)
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Printf("closing some client", err)
			for i, _ := range clients {
				if clients[i] == c {
					fmt.Println("removing client", len(clients))
					clients = append(clients[:i], clients[i+1:]...)
					fmt.Println("removed", len(clients))
					break
				}
			}
			c.Close()
			return
		}

		broadcast(netData, c)
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
