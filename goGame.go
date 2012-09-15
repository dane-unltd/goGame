package main

import (
	"github.com/dane-unltd/core"
	"os"
	"os/exec"
)

func main() {
	var host string = "localhost:33333"
	sp := false

	if len(os.Args) > 1 {
		if os.Args[1] == "-sp" {
			sp = true
			if len(os.Args) > 2 {
				host = os.Args[2]
			}
		} else {
			host = os.Args[1]
		}
	}

	if sp {
		cmd := exec.Command("srvSimple")
		cmd.Start()
		defer cmd.Wait()
	}

	sim := core.NewClient(os.Args[0], host)
	defer sim.Quit()

	sim.AddComp(false, 0, NewBoard(sp))
	sim.AddComp(true, 0, NewBrdDrawer())

	sim.Run()
}
