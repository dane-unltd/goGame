package main

import (
	//"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/dane-unltd/core"
	"time"
)

func main() {
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

	g.Init()
	g.Resize(int(screen.W), int(screen.H))

	sim := core.NewSim(g)

	input := new(core.Input)

	sim.AddComp("Input", 0, input)
	sim.AddComp("Board", 0, new(Board))
	sim.AddComp("BrdDrawer", 0, new(BrdDrawer))

	ticker := time.NewTicker(1e9 / 50 /*2 Hz*/)

	fmt.Println("entering loop")

loop:
	for {
		select {
		case <-ticker.C:
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
}
