package main

import "Server/internal/handler"

func main() {
	server := handler.NewServer()
	server.Start()
}
