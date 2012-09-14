package main

import "github.com/dane-unltd/core"

type BrdDrawer struct {
	board *Board
	g     core.Graphics
}

func (brdDrw *BrdDrawer) Id() string {
	return "BrdDrawer"
}

func NewBrdDrawer() *BrdDrawer {
	return &BrdDrawer{}
}

func (brdDrw *BrdDrawer) Swap() {
}

func (brdDrw *BrdDrawer) Init(g core.Graphics, sim *core.Sim, deps map[string]string) {
	brdDrw.g = g
	brdDrw.board = sim.Sys(deps["Board"]).(*Board)
}

func (brdDrw *BrdDrawer) Update() {
	brd := brdDrw.board.Board()
	var max float32 = BOARD * DIST
	for i := 0; i < BOARD; i++ {
		v := float32(i)*DIST + DIST/2
		brdDrw.g.DrawLine(0, v, max, v)
		brdDrw.g.DrawLine(v, 0, v, max)
	}

	for i := range brd {
		for j := range brd[i] {
			if brd[i][j] > 0 {
				brdDrw.draw(i, j, brd[i][j])
			}
		}
	}
}

func (brdDrw *BrdDrawer) draw(i, j int, player byte) {
	x := float32(i)*DIST + DIST/2
	y := float32(j)*DIST + DIST/2
	if player == 1 {
		brdDrw.g.DrawEntity(x, y, 0.0, "white")
	} else {
		brdDrw.g.DrawEntity(x, y, 0.0, "black")
	}
}
