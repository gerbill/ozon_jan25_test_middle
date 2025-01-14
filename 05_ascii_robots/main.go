package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var robotNames = []string{"A", "B"}
var wareHouse *WareHouse

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	Run(in, out)
}

func Run(in io.Reader, out io.Writer) {
	var dataSets, rows, columns int
	fmt.Fscan(in, &dataSets)
	for i := 1; i <= dataSets; i++ {
		fmt.Fscan(in, &rows, &columns)
		wareHouse = newWareHouse(rows, columns)
		for r := 0; r < rows; r++ {
			fmt.Fscan(in, &wareHouse.rows[r])
		}
		wareHouse.findRobots()
		wareHouse.pickDestination()
		wareHouse.move()
		wareHouse.report(out)
	}

}

type WareHouse struct {
	numRows    int
	numColumns int
	rows       []string
	columns    []string
	robots     []*Robot
	result     []string
}

func newWareHouse(numRows, numColumns int) *WareHouse {
	w := &WareHouse{
		numRows:    numRows,
		numColumns: numColumns,
		rows:       make([]string, numRows),
		columns:    make([]string, numColumns),
		result:     make([]string, numRows),
	}
	return w
}

func (w *WareHouse) move() {
	var counter int
	for !w.robots[0].arrived || !w.robots[1].arrived {
		for _, robot := range w.robots {
			robot.move()
		}
		counter++
		if counter > 1000000 {
			return
		}
	}
}

func (w *WareHouse) findRobots() {
	for y, row := range w.rows {
		for _, robotName := range robotNames {
			x := strings.Index(row, robotName)
			if x > -1 {
				r := newRobot(robotName, x, y)
				w.robots = append(w.robots, r)
			}
		}
	}
}

func (w *WareHouse) pickDestination() {
	A := w.robots[0]
	B := w.robots[1]
	if A.startY*A.startY+A.startX*A.startX > B.startY*B.startY+B.startX*B.startX {
		A.destX = w.numColumns - 1
		A.destY = w.numRows - 1
		B.destX = 0
		B.destY = 0
	} else {

		A.destX = 0
		A.destY = 0
		B.destX = w.numColumns - 1
		B.destY = w.numRows - 1
	}
}

func (w *WareHouse) makeColumns() {
	for _, row := range w.rows {
		for j := 0; j < len(row); j++ {
			w.columns[j] += string(row[j])
		}
	}
}

func (w *WareHouse) report(out io.Writer) {
	for _, row := range w.rows {
		fmt.Fprintln(out, row)
	}
}

type Robot struct {
	name                                             string
	step                                             string
	startX, startY, currentX, currentY, destX, destY int
	arrived                                          bool
}

func newRobot(name string, startX, startY int) *Robot {
	r := &Robot{
		name:     name,
		step:     strings.ToLower(name),
		startX:   startX,
		startY:   startY,
		currentX: startX,
		currentY: startY,
	}
	return r
}

func (r *Robot) move() {
	if r.currentX == r.destX && r.currentY == r.destY {
		r.arrived = true
		return
	}
	r.moveY()
}

func (r *Robot) moveX() {
	var move int
	if r.currentX < r.destX {
		move = 1
	} else if r.currentX > r.destX {
		move = -1
	}
	if canMove(r.currentX+move, r.currentY) {
		r.currentX += move
		replace(r.step, r.currentX, r.currentY)
	}
}

func (r *Robot) moveY() {
	var move int
	if r.currentY < r.destY {
		move = 1
	} else if r.currentY > r.destY {
		move = -1
	}
	if canMove(r.currentX, r.currentY+move) && move != 0 {
		r.currentY += move
		replace(r.step, r.currentX, r.currentY)
	} else {
		r.moveX()
		replace(r.step, r.currentX, r.currentY)
	}
}

func replace(replacement string, x, y int) {
	wareHouse.rows[y] = wareHouse.rows[y][:x] + replacement + wareHouse.rows[y][x+1:]
}

func canMove(x, y int) bool {
	if wareHouse.rows[y][x] != '.' {
		return false
	}
	return true
}
