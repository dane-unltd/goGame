package main

import (
	"github.com/dane-unltd/core"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	var host string = "localhost:33333"
	sp := false
	sz := 9

	l := len(os.Args)
	i := 1
	for i < l {
		a := os.Args[i]
		switch a {
		case "-sp":
			sp = true
		case "-p":
			i++
			sz, _ = strconv.Atoi(os.Args[i])
		default:
			host = a
		}
		i++
	}

	if sp {
		cmd := exec.Command("srvSimple")
		cmd.Start()
		defer cmd.Wait()
	}

	sim := core.NewClient(os.Args[0], host)
	defer sim.Quit()

	sim.AddComp(false, 0, NewBoard(sp, sz, 1.0))
	sim.AddComp(true, 0, NewBrdDrawer())

	sim.Run()
}
