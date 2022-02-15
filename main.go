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
	"net"
	"os"
)

// list of connected clients
var clients = []net.Conn{}

func addClient(c net.Conn) {
	fmt.Println("adding client:", c.RemoteAddr().String())
	clients = append(clients, c)
}

func removeClient(c net.Conn) {
	fmt.Println("removing client:", c.RemoteAddr().String())
	for i, _ := range clients {
		if clients[i] == c {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func broadcast(msg string, from net.Conn) {
	for _, c := range clients {
		if c != from {
			c.Write([]byte(string("< " + msg)))
		}
	}
}

func handleConnection(c net.Conn) {
	addClient(c)

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			removeClient(c)
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

	fmt.Println("Staring TCP Chat on port", PORT)
	for {
		// connect to new clients
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		// thread the process handler
		go handleConnection(c)
	}
}
