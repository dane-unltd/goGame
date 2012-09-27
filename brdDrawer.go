package main

import "github.com/dane-unltd/core"
import "github.com/dane-unltd/gina"

type BrdDrawer struct {
	board *Board
}

func (brdDrw *BrdDrawer) Id() string {
	return "BrdDrawer"
}

func NewBrdDrawer() *BrdDrawer {
	return &BrdDrawer{}
}

func (brdDrw *BrdDrawer) Swap() {
}

func (brdDrw *BrdDrawer) Init(sim core.Sim, res *core.ResMgr) {
	brdDrw.board = sim.Comp("Board").(*Board)
}

func (brdDrw *BrdDrawer) Update() {
	brd := brdDrw.board.Board()
	sz := brdDrw.board.Size()
	scale := 2000 / float32(sz)

	for i := 0; i < sz; i++ {
		v := -1000 + scale*(float32(i)+0.5)
		gina.DrawLine(-1000, v, 1000, v)
		gina.DrawLine(v, -1000, v, 1000)
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
	sz := brdDrw.board.Size()
	scale := 2000 / float32(sz)

	x := -1000 + scale*(float32(i)+0.5)
	y := -1000 + scale*(float32(j)+0.5)
	if player == 1 {
		gina.DrawEntity(x, y, 0.0, scale/2, "white")
	} else {
		gina.DrawEntity(x, y, 0.0, scale/2, "black")
	}
}
