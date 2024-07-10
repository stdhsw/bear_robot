package main

import (
	"os"
	"os/signal"
	"syscall"

	"example.com/bear/cmd/account"
	"example.com/bear/cmd/httpserver"
)

func main() {
	h := account.NewAccountHandler(os.Getenv("ACCOUNT_PATH"))
	server, err := httpserver.NewServer(h, os.Getenv("HTTP_PORT"))
	if err != nil {
		panic(err)
	}

	go server.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	server.Stop()
}
