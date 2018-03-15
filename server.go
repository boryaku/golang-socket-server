package main

import (
	"fmt"
	"net"
	"bufio"
	"awesomeProject/repos"
	"awesomeProject/handlers"
	"strings"
	"strconv"
)

const maxConnections = 5


/**
 * Bind tcp port 4000 and channel client connections to the handleConnection routine.
 */
func main() {
	server, err := net.Listen("tcp", ":4000")
	if server == nil {
		panic("couldn't start listening: " + err.Error())
	}

	//set up our number handler
	var values = make(chan string)
	go func() {
		numberHandler := handlers.NewNumberHandler(repos.NewNumberRepo())
		numberHandler.Save(values)
	}()

	//accept new socket connections
	connections := acceptClients(server)

	for {
		go stream(<-connections, values)
	}
}

/**
  * Accept client connections (max 5 clients)
 */
func acceptClients(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)

	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Printf("couldn't accept: " + err.Error())
				continue
			}
			i++
			if i > maxConnections {
				client.Close()
			} else {
				ch <- client
			}
		}
	}()

	return ch
}

/**
 * Read from the connection and write to the values stream.
 *
 * @todo Messages are separated by OS newline character
 */
func stream(conn net.Conn, values chan string) {
	connections := make([]net.Conn, maxConnections)
	connections = append(connections, conn)

	reader := bufio.NewReader(conn)

	for {
		message, _ := reader.ReadString('\n')

		if message == "terminate" {
			for _, connection := range connections {
				connection.Close()
			}
		}else if len(message) == 10 {
			//trim
			message = strings.TrimLeft(message, "0")
			message = strings.TrimRight(message, "\n")
			_, err := strconv.ParseUint(message, 0, 64)
			if err == nil {
				values <- message
			} else {
				conn.Close()
			}
		} else if len(message) > 0 && len(message) < 10{
			conn.Close()
		}
	}
}
