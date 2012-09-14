package main

import (
	//"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/dane-unltd/core"
	"os"
)

func main() {
	var host string

	if len(os.Args) > 1 {
		host = os.Args[1]
	} else {
		host = "localhost:33333"
	}
	sim := core.NewSim(os.Args[0], host)
	defer sim.Close()

	sim.AddComp(false, 0, NewBoard())
	sim.AddComp(true, 0, NewBrdDrawer())

	sim.Run()
}
