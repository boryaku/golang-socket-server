package main

import (
	"fmt"
	"net"
	"os"
	"math/rand"
	"time"
)

/**
 * Client to send messages to our server.
 * @todo control number of messages sent and handle sending terminate message
 */
func main() {
	rand.Seed(time.Now().Unix())
	fmt.Println("start time =", time.Now())

	conn, _ := net.Dial("tcp", ":4000")

	for n := 0; n < 20000000; n++ {
		_, err := fmt.Fprintf(conn,  generateValue()+ "\n")
		if err != nil {
			fmt.Println("having problems writing to the socket, going to disconnect")
			n = 20000000
			conn.Close()
		}
	}

	fmt.Println("end time =", time.Now())

	os.Exit(0)
}

/**
 * Generate a random sting of 9 digits.
 */
func generateValue() string {
	var numbers = []rune("0123456789")

	b := make([]rune, 9)
	for i := range b {
		b[i] = numbers[rand.Intn(len(numbers))]
	}

	return string(b)
}

