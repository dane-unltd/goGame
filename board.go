package main

import "github.com/dane-unltd/core"

import "fmt"

type Point struct {
	x, y int
}

type boardState struct {
	brd, nbrs [][]byte
	pass, fin bool
	plr       uint32
	points    [2]int
}

type Board struct {
	curr, next *boardState
	helper     [][]bool
	cmdSrc     core.CmdSrc
	sp         bool
	sz         int
}

func NewBoardState(size int) *boardState {
	state := &boardState{
		make([][]byte, size),
		make([][]byte, size),
		false, false,
		1,
		[2]int{0, 0},
	}

	for i := range state.brd {
		state.brd[i] = make([]byte, size)
		state.nbrs[i] = make([]byte, size)
	}

	for i := 0; i < size; i++ {
		state.nbrs[i][0] += 1
		state.nbrs[i][size-1] += 1
		state.nbrs[0][i] += 1
		state.nbrs[size-1][i] += 1
	}
	return state
}

func NewBoard(sp bool, size int) *Board {
	board := &Board{
		NewBoardState(size), NewBoardState(size),
		make([][]bool, size),
		nil,
		sp,
		size,
	}

	for i := range board.helper {
		board.helper[i] = make([]bool, board.sz)
	}

	return board
}

func (board *Board) Board() [][]byte {
	return board.curr.brd
}

func (board *Board) Size() int {
	return board.sz
}

func (board *Board) Result() (bool, [2]int) {
	return board.curr.fin, board.curr.points
}

func (board *Board) Id() string {
	return "Board"
}

func (board *Board) Init(sim core.Sim, res *core.ResMgr, deps map[string]string) {
	board.cmdSrc = sim.GetComp(deps["CmdSrc"]).(core.CmdSrc)
}

func (board *Board) Swap() {
	temp := board.curr
	board.curr = board.next
	board.next = temp
}

func (board *Board) Update() {
	copy(board.next.brd, board.curr.brd)
	copy(board.next.nbrs, board.curr.nbrs)

	board.next.plr = board.curr.plr
	board.next.fin = board.curr.fin
	board.next.pass = board.curr.pass
	board.next.points = board.curr.points

	if board.curr.fin {
		return
	}

	var pId uint32
	if board.sp {
		pId = 1
	} else {
		pId = board.curr.plr
	}

	if board.cmdSrc.Active(pId, "pass") {
		board.next.plr = board.curr.plr%2 + 1
		board.next.pass = true
		if board.curr.pass == true {
			board.next.fin = true
			board.calcScore()
			fmt.Println("finished:", board.next.points)
		}

		return
	}

	if board.cmdSrc.Active(pId, "place") {
		X, Y := board.cmdSrc.Point(pId)
		x := ((X + 1000) * board.sz / 2000)
		y := ((Y + 1000) * board.sz / 2000)

		if x >= 0 && x < board.sz && y >= 0 && y < board.sz {
			if board.place(x, y) {
				board.next.plr = board.curr.plr%2 + 1
			}
		}
	}
}

func (board *Board) place(x, y int) bool {
	if board.next.brd[x][y] == 0 {
		plr := board.curr.plr
		nbrs := board.neighbours(x, y)

		board.next.brd[x][y] = byte(plr)
		for _, nbr := range nbrs {
			board.next.nbrs[nbr.x][nbr.y]++
		}

		capt := false
		for _, nbr := range nbrs {
			if board.next.brd[nbr.x][nbr.y] & ^byte(plr) > 0 {
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
	plr := board.next.brd[x][y]
	chain = make([]Point, 1)
	chain[0] = Point{x, y}
	board.helper[x][y] = true
	cnt = 0
	for cnt < len(chain) {
		pt := chain[cnt]
		if board.next.nbrs[pt.x][pt.y] < 4 {
			cnt = 0
			break
		}
		nbrs := board.neighbours(pt.x, pt.y)

		for _, nbr := range nbrs {
			if board.next.brd[nbr.x][nbr.y] == plr {
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
	if x < board.sz-1 {
		nbrs = append(nbrs, Point{x + 1, y})
	}
	if y < board.sz-1 {
		nbrs = append(nbrs, Point{x, y + 1})
	}
	return nbrs
}

func (board *Board) remove(chain []Point) {
	for _, pt := range chain {
		nbrs := board.neighbours(pt.x, pt.y)
		board.next.brd[pt.x][pt.y] = 0
		for _, nbr := range nbrs {
			board.next.nbrs[nbr.x][nbr.y]--
		}
	}
}

func (board *Board) calcScore() {

	for x := 0; x < board.sz; x++ {
		for y := 0; y < board.sz; y++ {
			cnt, plr := board.checkTerr(x, y)
			if plr > 0 {
				board.next.points[plr-1] += cnt
			}
		}
	}
}

func (board *Board) checkTerr(x, y int) (cnt int, plr byte) {
	plr = board.curr.brd[x][y]
	if plr > 0 {
		return 1, plr
	}
	if board.helper[x][y] == true {
		return 0, 0
	}

	chain := make([]Point, 1)
	chain[0] = Point{x, y}
	board.helper[x][y] = true
	cnt = 0
	neutral := false
	for cnt < len(chain) {
		pt := chain[cnt]

		nbrs := board.neighbours(pt.x, pt.y)

		for _, nbr := range nbrs {
			curr := board.curr.brd[nbr.x][nbr.y]
			if curr == 0 {
				if board.helper[nbr.x][nbr.y] == false {
					chain = append(chain, nbr)
					board.helper[nbr.x][nbr.y] = true
				}
			}
			if !neutral {
				if curr > 0 {
					if plr == 0 {
						plr = curr
					} else if plr != curr {
						neutral = true
						plr = 0
					}
				}
			}
		}

		cnt++
	}
	return
}
