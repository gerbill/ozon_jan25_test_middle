package main

import (
	"bufio"
	"bytes"
	"cmp"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	Run(in, out)
}

func Run(in *bufio.Reader, out *bufio.Writer) {
	lineCounter := 0
	ob := OrderBatch{}
	for {
		lineCounter++
		line, err := in.ReadString(byte('\n'))
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read line %d. Error: %+v", lineCounter, err)
		}
		if line == "" {

			break
		}
		if lineCounter == 1 {
			continue
		}
		// l := lineCounter - 2
		line = line[:len(line)-1]
		if !ob.isComplete {
			ob.fill(line)
		} else {
			out.WriteString(ob.report() + "\n")
			ob = OrderBatch{}
			ob.fill(line)
		}
	}
	if ob.isComplete {
		out.WriteString(ob.report() + "\n")
	}
}

var sortedTrucks []*Truck
var maxWindowDuration = 0
var lastAvailableTruckId = 0

type OrderBatch struct {
	ordersNum  int
	orders     []*Order
	orders2    []*Order
	trucksNum  int
	trucks     []*Truck
	isComplete bool
}

type Truck struct {
	id       int
	start    int
	end      int
	capacity int
	orders   int
}

type Order struct {
	id      int
	arrival int
	truckId int
}

func (ob *OrderBatch) fill(line string) {
	switch {
	case ob.ordersNum == 0:
		ob.ordersNum = aToI(line, "orders num")
		ob.orders = make([]*Order, 0, ob.ordersNum)
	case len(ob.orders) != ob.ordersNum:
		arrivals := strings.Fields(line)
		for i, arrival := range arrivals {
			order := newOrder(i+1, arrival)
			ob.orders = append(ob.orders, order)
			ob.orders2 = append(ob.orders2, order)
		}
	case ob.trucksNum == 0:
		ob.trucksNum = aToI(line, "trucks num")
		ob.trucks = make([]*Truck, 0, ob.trucksNum)
	case len(ob.trucks) != ob.trucksNum:
		truckProps := strings.Fields(line)
		if len(truckProps) != 3 {
			log.Fatalf("failed to parse truckProps on line %s. Expected 3, got %d", line, len(truckProps))
		}
		ob.trucks = append(ob.trucks, newTruck(len(ob.trucks)+1, truckProps[0], truckProps[1], truckProps[2]))
	}
	if len(ob.trucks) == ob.trucksNum && len(ob.trucks) != 0 {
		ob.isComplete = true
		sortedTrucks = []*Truck{}
		maxWindowDuration = 0
		lastAvailableTruckId = 0
	}
}

func aToI(a string, message string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		log.Fatalf("failed to convert string '%s' to int. Error: %+v. Message: %s", a, err, message)
	}
	return i
}

func (ob *OrderBatch) report() string {
	var report bytes.Buffer
	sortedTrucks = ob.sortedTrucks()
	sortedOrders := ob.sortedOrders()
	for _, order := range sortedOrders {
		order.findTruck()
	}
	for _, order := range ob.orders2 {
		report.WriteString(strconv.Itoa(order.truckId) + " ")
	}
	return report.String()
}

func (ob *OrderBatch) sortedTrucks() []*Truck {
	slices.SortFunc(ob.trucks, func(a, b *Truck) int {
		return cmp.Or(
			cmp.Compare(a.start, b.start),
			cmp.Compare(a.id, b.id),
		)
	})
	return ob.trucks
}

func (ob *OrderBatch) sortedOrders() []*Order {
	slices.SortFunc(ob.orders, func(a, b *Order) int {
		return cmp.Or(
			cmp.Compare(a.arrival, b.arrival),
		)
	})
	return ob.orders
}

// Order //////////////////////////////////////////////////
func newOrder(id int, arrival string) *Order {
	return &Order{
		id:      id,
		arrival: aToI(arrival, "order arrival"),
	}
}

func (o *Order) findTruck() {
	offset := lastAvailableTruckId
	for i, truck := range sortedTrucks[lastAvailableTruckId:] {
		if truck.start > o.arrival {
			break
		}
		if truck.canTakeOrder(o.arrival) {
			o.truckId = truck.id
			return
		}
		if truck.start+maxWindowDuration < o.arrival {
			lastAvailableTruckId = i + offset
		}
	}
	o.truckId = -1
}

func (o *Order) findMinTruck() int {
	min := (len(sortedTrucks) - 1) / 2
	for {
		switch {
		case sortedTrucks[min].start <= o.arrival:
			return sortedTrucks[min].id
		}
	}
}

// Truck //////////////////////////////////////////////////
func newTruck(id int, start, end, capacity string) *Truck {
	t := &Truck{
		id:       id,
		start:    aToI(start, "truck start"),
		end:      aToI(end, "truck end"),
		capacity: aToI(capacity, "truck capacity"),
	}
	window := t.end - t.start
	if window > maxWindowDuration {
		maxWindowDuration = window
	}
	return t
}

func (t *Truck) canTakeOrder(arrival int) bool {
	if arrival < t.start {
		return false
	}
	if arrival > t.end {
		return false
	}
	if t.orders >= t.capacity {
		return false
	}
	t.orders++
	return true
}

func (t *Truck) freeSpace() int {
	return t.capacity - t.orders
}
