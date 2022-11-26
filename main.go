package main

import (
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// get arguments
	args := os.Args
	if len(args) < 2 {
		panic("arguments at least one needed")
	}

	// creating folder
	in := strings.Split(args[1], "/")
	dir := strings.Join(in[:len(in)-1], "/")

	// if err := exec.Command("rm", "-rf", dir).Run(); err != nil {
	// 	panic(err)
	// }

	if err := exec.Command("mkdir", "-p", dir).Run(); err != nil {
		panic(err)
	}

	// process binding and listening
	listener, err := net.Listen("unix", args[1])
	if err != nil {
		log.Fatal("listen error:", err)
	}
	log.Println("server ready to accepting connection")

	grace := make(chan os.Signal, 1)
	signal.Notify(grace, syscall.SIGINT, syscall.SIGTERM)

	defer listener.Close()
	go echoServer(listener)

	// grace shutting down
	<-grace

	log.Println("removing socket file")
	exec.Command("rm", "-f", args[1]).Run()
}
