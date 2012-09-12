package main

import "github.com/dane-unltd/core"

type BrdDrawer struct {
	board *Board
	g     core.Graphics
}

func (this *BrdDrawer) Id() string {
	return "BrdDrawer"
}

func NewBrdDrawer(sim *core.Sim, g core.Graphics) *BrdDrawer {
	return &BrdDrawer{
		sim.Sys("Board").(*Board),
		g,
	}
}

func (this *BrdDrawer) Swap() {
}

func (this *BrdDrawer) Update() {
	brd := this.board.Board()
	var max float32 = BOARD * DIST
	for i := 0; i < BOARD; i++ {
		v := float32(i)*DIST + DIST/2
		this.g.DrawLine(0, v, max, v)
		this.g.DrawLine(v, 0, v, max)
	}

	for i := range brd {
		for j := range brd[i] {
			if brd[i][j] > 0 {
				this.draw(i, j, brd[i][j])
			}
		}
	}
}

func (this *BrdDrawer) draw(i, j int, player byte) {
	x := float32(i)*DIST + DIST/2
	y := float32(j)*DIST + DIST/2
	if player == 1 {
		this.g.DrawEntity(x, y, 0.0, "white")
	} else {
		this.g.DrawEntity(x, y, 0.0, "black")
	}
}
