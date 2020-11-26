package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("must be `watch <path> <command>`")
	}

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	changes, close := WatchFolder(os.Args[1])
	defer close()

	par := &Par{
		kill:     make(chan bool),
		commands: os.Args[2:],
	}

	go par.run()

	for {
		select {
		case <-interrupt:
			par.kill <- true
			os.Exit(0)
		case <-changes:
			par.kill <- true
			go par.run()
		}
	}
}
