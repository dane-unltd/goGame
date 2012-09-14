package main

import (
	//"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/dane-unltd/core"
	"net"
	"os"
	"runtime"
	"time"
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

	deps := map[string]map[string]string{
		"Input":     map[string]string{"TimeSrc": "CltRx"},
		"Board":     map[string]string{"CmdSrc": "CltRx"},
		"CltTx":     map[string]string{"CmdSrc": "Input"},
		"BrdDrawer": map[string]string{"Board": "Board"},
	}
	sim.AddComp(false, 0, rcv)
	sim.AddComp(true, 0, input)
	sim.AddComp(true, 0, core.NewCltTx(conn))
	sim.AddComp(false, 0, NewBoard())
	sim.AddComp(true, 0, NewBrdDrawer(g))

	sim.UpdateDeps(deps)

	fmt.Println("entering loop")

	rcv.RxLocalId()

	clk := time.Tick(rcv.FrameNS())

	go sim.RunAsync(rcv)

loop:
	for {
		<-clk

		g.PreFrame()
		sim.SyncUpdate()

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
