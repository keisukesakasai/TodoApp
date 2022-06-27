// main.go

package main

import (
	"TodoApp/app/controllers"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Process Shutdown...")
		os.Exit(1)
	}()
}

func main() {
	controllers.StartMainServer()
}
