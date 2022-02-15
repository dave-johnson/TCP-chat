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

type Client struct {
	Conn net.Conn
	Name string
}

// list of connected clients
var clients = []Client{}

func addClient(n string, c net.Conn) {
	fmt.Printf("adding client {%s} from %s\n", n, c.RemoteAddr().String())
	client := Client{
		Conn: c,
		Name: n,
	}
	clients = append(clients, client)
}

func removeClient(c net.Conn) {
	fmt.Println("removing client:", c.RemoteAddr().String())
	for i, _ := range clients {
		if clients[i].Conn.RemoteAddr() == c.RemoteAddr() {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func broadcast(msg string, from net.Conn) {
	for _, c := range clients {
		if c.Conn != from {
			c.Conn.Write([]byte(string("<" + c.Name + "> " + msg)))
		}
	}
}

func handleConnection(c net.Conn) {
	addClient(c.RemoteAddr().String(), c)

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
