package main

import "github.com/dane-unltd/core"

//import "fmt"

const (
	BOARD = 7
	DIST  = 50
)

type Point struct {
	x, y int
}
type Board struct {
	currNbrs, nextNbrs   [][]byte
	currBoard, nextBoard [][]byte

	helper                 [][]bool
	currPlayer, nextPlayer byte
	input                  *core.Input
}

func (this *Board) Board() [][]byte {
	return this.currBoard
}

func (this *Board) Id() string {
	return "Board"
}

func (this *Board) Init(g core.Graphics, sim *core.Sim) {
	this.currPlayer = 1
	this.input = sim.Sys("Input").(*core.Input)

	this.currBoard = make([][]byte, BOARD)
	this.nextBoard = make([][]byte, BOARD)
	this.currNbrs = make([][]byte, BOARD)
	this.nextNbrs = make([][]byte, BOARD)
	this.helper = make([][]bool, BOARD)
	for i := range this.currBoard {
		this.currBoard[i] = make([]byte, BOARD)
		this.nextBoard[i] = make([]byte, BOARD)
		this.currNbrs[i] = make([]byte, BOARD)
		this.nextNbrs[i] = make([]byte, BOARD)
		this.helper[i] = make([]bool, BOARD)
	}

	for i := 0; i < BOARD; i++ {
		this.currNbrs[i][0] += 1
		this.currNbrs[i][BOARD-1] += 1
		this.currNbrs[0][i] += 1
		this.currNbrs[BOARD-1][i] += 1
	}
}

func (this *Board) Swap() {
	temp := this.currBoard
	this.currBoard = this.nextBoard
	this.nextBoard = temp

	temp = this.currNbrs
	this.currNbrs = this.nextNbrs
	this.nextNbrs = temp

	this.currPlayer = this.nextPlayer
}

func (this *Board) Update() {
	copy(this.nextBoard, this.currBoard)
	copy(this.nextNbrs, this.currNbrs)

	this.nextPlayer = this.currPlayer
	cmd := this.input.Cmd()

	if cmd.Actions&core.ACTION1 > 0 {
		x := cmd.X / DIST
		y := cmd.Y / DIST
		if x >= 0 && x < BOARD && y >= 0 && y < BOARD {
			if this.place(x, y) {
				this.nextPlayer = this.currPlayer%2 + 1
			}
		}
	}
}

func (this *Board) place(x, y int) bool {
	if this.nextBoard[x][y] == 0 {
		plr := this.currPlayer
		nbrs := this.neighbours(x, y)

		this.nextBoard[x][y] = plr
		for _, nbr := range nbrs {
			this.nextNbrs[nbr.x][nbr.y]++
		}

		capt := false
		for _, nbr := range nbrs {
			if this.nextBoard[nbr.x][nbr.y] & ^plr > 0 {
				cnt, chain := this.capture(nbr.x, nbr.y)
				if cnt > 0 {
					capt = true
					this.remove(chain)
				}
			}
		}

		//handle suicide
		if !capt {
			if cnt, _ := this.capture(x, y); cnt > 0 {
				this.remove([]Point{{x, y}})
				return false
			}
		}

		return true
	}
	return false
}

func (this *Board) capture(x, y int) (cnt int, chain []Point) {
	plr := this.nextBoard[x][y]
	chain = make([]Point, 1)
	chain[0] = Point{x, y}
	this.helper[x][y] = true
	cnt = 0
	for cnt < len(chain) {
		pt := chain[cnt]
		if this.nextNbrs[pt.x][pt.y] < 4 {
			cnt = 0
			break
		}
		nbrs := this.neighbours(pt.x, pt.y)

		for _, nbr := range nbrs {
			if this.nextBoard[nbr.x][nbr.y] == plr {
				if this.helper[nbr.x][nbr.y] == false {
					chain = append(chain, nbr)
					this.helper[nbr.x][nbr.y] = true
				}
			}
		}

		cnt++
	}

	for _, pt := range chain {
		this.helper[pt.x][pt.y] = false
	}
	return
}

func (this *Board) neighbours(x, y int) []Point {
	nbrs := make([]Point, 0)
	if x > 0 {
		nbrs = append(nbrs, Point{x - 1, y})
	}
	if y > 0 {
		nbrs = append(nbrs, Point{x, y - 1})
	}
	if x < BOARD-1 {
		nbrs = append(nbrs, Point{x + 1, y})
	}
	if y < BOARD-1 {
		nbrs = append(nbrs, Point{x, y + 1})
	}
	return nbrs
}

func (this *Board) remove(chain []Point) {
	for _, pt := range chain {
		nbrs := this.neighbours(pt.x, pt.y)
		this.nextBoard[pt.x][pt.y] = 0
		for _, nbr := range nbrs {
			this.nextNbrs[nbr.x][nbr.y]--
		}
	}
}
