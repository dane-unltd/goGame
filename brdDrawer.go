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

func (brdDrw *BrdDrawer) Init(g core.Graphics, sim core.Sim, deps map[string]string) {
	brdDrw.g = g
	brdDrw.board = sim.GetComp(deps["Board"]).(*Board)
	sz, scale := brdDrw.board.Params()
	g.Resize(int(50.0*float32(sz)*scale), int(50.0*float32(sz)*scale))
}

func (brdDrw *BrdDrawer) Update() {
	brd := brdDrw.board.Board()
	sz, scale := brdDrw.board.Params()
	var max float32 = 50.0 * float32(sz) * scale
	for i := 0; i < sz; i++ {
		v := 50.0 * scale * (float32(i) + 0.5)
		brdDrw.g.DrawLine(0, v, max, v)
		brdDrw.g.DrawLine(v, 0, v, max)
	}

	for i := range brd {
		for j := range brd[i] {
			if brd[i][j] > 0 {
				brdDrw.draw(i, j, scale, brd[i][j])
			}
		}
	}

	fin, points := brdDrw.board.Result()
	if fin {
		brdDrw.g.DrawHS(points[:])
	}
}

func (brdDrw *BrdDrawer) draw(i, j int, scale float32, player byte) {
	x := 50.0 * scale * (float32(i) + 0.5)
	y := 50.0 * scale * (float32(j) + 0.5)
	if player == 1 {
		brdDrw.g.DrawEntity(x, y, 0.0, scale, "white")
	} else {
		brdDrw.g.DrawEntity(x, y, 0.0, scale, "black")
	}
}
