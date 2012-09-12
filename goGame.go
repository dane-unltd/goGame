package main

import (
	//"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/dane-unltd/core"
	"net"
	"os"
	"runtime"
)

func main() {
	runtime.LockOSThread()
	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		panic(sdl.GetError())
	}

	defer sdl.Quit()

	screen := sdl.SetVideoMode(BOARD*DIST, BOARD*DIST, 32, sdl.OPENGL)
	if screen == nil {
		panic(sdl.GetError())
	}
	sdl.WM_SetCaption("Engine", "")

	g := core.NewGraphics()

	g.Resize(int(screen.W), int(screen.H))

	var host string

	if len(os.Args) > 1 {
		host = os.Args[1]
	} else {
		host = "localhost:33333"
	}
	conn, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	rcv := core.NewCltRx(conn)

	sim := core.NewSim(g)

	input := core.NewInput()

	sim.AddComp(0, rcv)
	sim.AddComp(0, input)
	sim.AddComp(0, core.NewCltTx(conn, sim))
	sim.AddComp(0, NewBoard(sim))
	sim.AddComp(0, NewBrdDrawer(sim, g))

	fmt.Println("entering loop")

loop:
	for {
		err = rcv.Receive()
		if err != nil {
			fmt.Println(err)
			break
		}

		g.PreFrame()
		sim.Update()

		if input.Quit() {
			break loop
		}
		if r, w, h := input.Resize(); r {
			screen = sdl.SetVideoMode(w, h, 32, sdl.OPENGL)
			g.Resize(w, h)
		}

		g.PostFrame()

		sdl.GL_SwapBuffers()
	}
}
