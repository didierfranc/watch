package main

import (
	"log"
	"os"
	"os/signal"

	c "github.com/didierfranc/watch/pkg/command"
	w "github.com/didierfranc/watch/pkg/watch"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("must be `watch <path> <command>`")
	}

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	changes, close := w.WatchFolder(os.Args[1])
	defer close()

	par := &c.Par{
		Kill:     make(chan bool),
		Commands: os.Args[2:],
	}

	go par.Run()

	for {
		select {
		case <-interrupt:
			par.Kill <- true
			os.Exit(0)
		case <-changes:
			par.Kill <- true
			go par.Run()
		}
	}
}
