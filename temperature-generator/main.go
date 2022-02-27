package main

import (
	"fmt"
	"github.com/jponge/playground-go-microservices/temperature-generator/cmd"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
	if !cmd.RunStarted() {
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	fmt.Println("👋 Bye!")
}
