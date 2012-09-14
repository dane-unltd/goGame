package main

import "github.com/dane-unltd/core"

import "fmt"

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
	currPlayer, nextPlayer uint32
	cmdSrc                 core.CmdSrc
}

func (board *Board) Board() [][]byte {
	return board.currBoard
}

func (board *Board) Id() string {
	return "Board"
}

func NewBoard() *Board {
	board := &Board{}

	board.currPlayer = 1

	board.currBoard = make([][]byte, BOARD)
	board.nextBoard = make([][]byte, BOARD)
	board.currNbrs = make([][]byte, BOARD)
	board.nextNbrs = make([][]byte, BOARD)
	board.helper = make([][]bool, BOARD)
	for i := range board.currBoard {
		board.currBoard[i] = make([]byte, BOARD)
		board.nextBoard[i] = make([]byte, BOARD)
		board.currNbrs[i] = make([]byte, BOARD)
		board.nextNbrs[i] = make([]byte, BOARD)
		board.helper[i] = make([]bool, BOARD)
	}

	for i := 0; i < BOARD; i++ {
		board.currNbrs[i][0] += 1
		board.currNbrs[i][BOARD-1] += 1
		board.currNbrs[0][i] += 1
		board.currNbrs[BOARD-1][i] += 1
	}
	return board
}

func (board *Board) Init(g core.Graphics, sim *core.Sim, deps map[string]string) {
	board.cmdSrc = sim.Sys(deps["CmdSrc"]).(core.CmdSrc)
}

func (board *Board) Swap() {
	temp := board.currBoard
	board.currBoard = board.nextBoard
	board.nextBoard = temp

	temp = board.currNbrs
	board.currNbrs = board.nextNbrs
	board.nextNbrs = temp

	board.currPlayer = board.nextPlayer
}

func (board *Board) Update() {
	copy(board.nextBoard, board.currBoard)
	copy(board.nextNbrs, board.currNbrs)

	board.nextPlayer = board.currPlayer
	cmd := board.cmdSrc.Cmd(board.currPlayer)

	fmt.Println("Board:", cmd)

	if cmd.Actions&core.ACTION1 > 0 {
		x := cmd.X / DIST
		y := cmd.Y / DIST
		if x >= 0 && x < BOARD && y >= 0 && y < BOARD {
			if board.place(x, y) {
				board.nextPlayer = board.currPlayer%2 + 1
			}
		}
	}
}

func (board *Board) place(x, y int) bool {
	if board.nextBoard[x][y] == 0 {
		plr := board.currPlayer
		nbrs := board.neighbours(x, y)

		board.nextBoard[x][y] = byte(plr)
		for _, nbr := range nbrs {
			board.nextNbrs[nbr.x][nbr.y]++
		}

		capt := false
		for _, nbr := range nbrs {
			if board.nextBoard[nbr.x][nbr.y] & ^byte(plr) > 0 {
				cnt, chain := board.capture(nbr.x, nbr.y)
				if cnt > 0 {
					capt = true
					board.remove(chain)
				}
			}
		}

		//handle suicide
		if !capt {
			if cnt, _ := board.capture(x, y); cnt > 0 {
				board.remove([]Point{{x, y}})
				return false
			}
		}

		return true
	}
	return false
}

func (board *Board) capture(x, y int) (cnt int, chain []Point) {
	plr := board.nextBoard[x][y]
	chain = make([]Point, 1)
	chain[0] = Point{x, y}
	board.helper[x][y] = true
	cnt = 0
	for cnt < len(chain) {
		pt := chain[cnt]
		if board.nextNbrs[pt.x][pt.y] < 4 {
			cnt = 0
			break
		}
		nbrs := board.neighbours(pt.x, pt.y)

		for _, nbr := range nbrs {
			if board.nextBoard[nbr.x][nbr.y] == plr {
				if board.helper[nbr.x][nbr.y] == false {
					chain = append(chain, nbr)
					board.helper[nbr.x][nbr.y] = true
				}
			}
		}

		cnt++
	}

	for _, pt := range chain {
		board.helper[pt.x][pt.y] = false
	}
	return
}

func (board *Board) neighbours(x, y int) []Point {
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

func (board *Board) remove(chain []Point) {
	for _, pt := range chain {
		nbrs := board.neighbours(pt.x, pt.y)
		board.nextBoard[pt.x][pt.y] = 0
		for _, nbr := range nbrs {
			board.nextNbrs[nbr.x][nbr.y]--
		}
	}
}
