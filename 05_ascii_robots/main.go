package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var robotNames = []string{"A", "B"}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	Run(in, out)
}

func Run(in io.Reader, out io.Writer) {
	var dataSets, rows, columns int
	fmt.Fscan(in, &dataSets)
	fmt.Println(dataSets)
	for i := 1; i <= dataSets; i++ {
		fmt.Fscan(in, &rows, &columns)
		w := newWareHouse(rows, columns)
		for r := 0; r < rows; r++ {
			fmt.Fscan(in, &w.rows[r])
		}
		w.findRobots()
		w.pickDestination()
		w.move()
		w.report(out)
	}

}

type WareHouse struct {
	numRows    int
	numColumns int
	rows       []string
	robots     []*Robot
	result     []string
}

func newWareHouse(numRows, numColumns int) *WareHouse {
	return &WareHouse{
		numRows:    numRows,
		numColumns: numColumns,
		rows:       make([]string, numRows),
		result:     make([]string, numRows),
	}
}

func (w *WareHouse) move() {
	for i := 0; i < w.numRows; i++ {
		for _, r := range w.robots {
			result := r.move(i, w.rows[i])
			if result != "" {
				w.result[i] = result
			}
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
	if A.startY < B.startX {
		A.destX = 0
		A.destY = 0
		A.movingUp = true
		B.destX = w.numColumns - 1
		B.destY = w.numRows - 1
		B.movingUp = false
	} else {
		A.destX = w.numColumns - 1
		A.destY = w.numRows - 1
		A.movingUp = false
		B.destX = 0
		B.destY = 0
		B.movingUp = true
	}
	fmt.Printf("robot %+v\n", A)
	fmt.Printf("robot %+v\n", B)
	fmt.Println("==========")
}

func (w *WareHouse) report(out io.Writer) {
	for _, row := range w.result {
		fmt.Fprintln(out, row)
	}
}

type Robot struct {
	name                                             string
	step                                             string
	movingUp                                         bool
	startX, startY, currentX, currentY, destX, destY int
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

func (r *Robot) move(rowId int, row string) string {
	l := len(row)
	if rowId == r.startY && r.startX == r.destX && r.startY == r.destY {
		return row
	}
	if r.movingUp && rowId > r.startY {
		return ""
	}
	if rowId == r.destY {
		if r.movingUp {
			return strings.Repeat(r.step, r.startX+1) + strings.Repeat(".", l-r.startX)
		}
		return strings.Repeat(".", r.startX) + strings.Repeat(r.step, r.startX)
	}
	return row[:r.startX+1] + r.step + row[r.startX+1:]
}
