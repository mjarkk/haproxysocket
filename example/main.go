package main

import (
	"fmt"

	"github.com/mjarkk/haproxysocket"
)

func main() {
	// Create a instace of haproxy
	// Make sure to change the haproxy/haproxy.sock to where your haproxy sock file is
	h := haproxysocket.New("unix", "../testing/haproxy/haproxy.sock")

	// Get the sessions
	sessions, err := h.ShowSess()
	if err != nil {
		panic(err)
	}
	fmt.Println("Sessions:")
	for _, session := range sessions {
		fmt.Println("ID:", session.ID)
	}

	// Set a server to maintenance mode
	s := h.Server("test-backend", "serv1")
	err = s.State("maint")
	if err != nil {
		panic(err)
	}
}
