// Should be able to receive and forward messages among clients
// > yo
// < whats up?
// > Nothing much!!
// < same here, hbu?
// > same ol' same

package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	// chat commands
	JOIN   = "join"
	WHOAMI = "whoami"
	HELP   = "help"

	// admin is a special user
	ADMIN = "admin"
	PWD   = "password" // this will come from an environemnt variable later
)

type Client struct {
	Conn net.Conn
	Name string
}

// list of connected clients
var clients = []Client{}

func join(d []string, c net.Conn) (string, error) {
	if len(d) == 1 {
		c.Write([]byte("You need to enter your username\n"))
		return "", errors.New("invalid join statement")
	}

	if d[1] == ADMIN {
		if len(d) == 2 {
			c.Write([]byte("You need to enter a password\n"))
			return "", errors.New("invalid join statement")
		}
		if d[2] != PWD {
			c.Write([]byte("You entered the wrong password\n"))
			return "", errors.New("invalid join statement")
		}
		// warning message
		if len(d) > 3 {
			c.Write([]byte("Values after password will be ignored\n"))
		}
	}

	addClient(d[1], c)
	return d[1], nil
}

func addClient(n string, c net.Conn) {
	fmt.Printf("adding client {%s} from %s\n", n, c.RemoteAddr().String())
	client := Client{
		Conn: c,
		Name: n,
	}
	clients = append(clients, client)
}

func removeClient(c net.Conn) {
	for i, _ := range clients {
		if clients[i].Conn.RemoteAddr() == c.RemoteAddr() {
			fmt.Println("removing client:", clients[i].Name)
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func broadcast(msg string, from string) {
	fmt.Println("------ broadcast from ", from)
	for _, c := range clients {
		if c.Name != from {
			c.Conn.Write([]byte(string("<" + from + "> " + msg)))
		}
	}
}

func handleConnection(c net.Conn) {
	var from = string("")

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			removeClient(c)
			c.Close()
			return
		}

		// don't broadcast empty messages
		if netData == "\n" {
			continue
		}

		d := strings.Split(strings.TrimRight(netData, "\n"), " ")
		switch {
		case d[0] == JOIN:
			from, _ = join(d, c)
		case from == "":
			c.Write([]byte("You must JOIN and enter your user name before posting messages.\n"))
		default:
			broadcast(netData, from)
		}

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
