package main

import (
	"bufio"
	"bytes"
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	Run(in, out)
}

var sortedTrucks []*Truck

type OrderBatch struct {
	ordersNum  int
	orders     []*Order
	orders2    []*Order
	trucksNum  int
	trucks     []*Truck
	isComplete bool
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
		// ob.printTrucks()
	}
}

func (ob *OrderBatch) report() string {
	var report bytes.Buffer
	ts := time.Now().UnixNano()
	sortedTrucks = ob.sortedTrucks()
	fmt.Println("sort trucks", time.Now().UnixNano()-ts)
	ts = time.Now().UnixNano()
	sortedOrders := ob.sortedOrders()
	fmt.Println("sort orders", time.Now().UnixNano()-ts)
	ts = time.Now().UnixNano()
	printTrucksStart()
	printOrders(sortedOrders)
	for _, order := range sortedOrders {
		order.truckId = order.findTruck()
	}
	fmt.Println("find trucks for orders", time.Now().UnixNano()-ts)
	ts = time.Now().UnixNano()
	for _, order := range ob.orders2 {
		report.WriteString(strconv.Itoa(order.truckId) + " ")
	}
	fmt.Println("add to report", time.Now().UnixNano()-ts)
	return report.String()
}

func (ob *OrderBatch) sortedTrucks() []*Truck {
	slices.SortFunc(ob.trucks, func(a, b *Truck) int {
		return cmp.Or(
			cmp.Compare(a.start, b.start),
			// cmp.Compare(b.freeSpace(), a.freeSpace()),
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

type Order struct {
	id      int
	arrival int
	truckId int
}

func (o *Order) findTruck() int {
	for _, truck := range sortedTrucks {
		if truck.canTakeOrder(o.arrival) {
			return truck.id
		}
	}
	return -1
}

func newOrder(id int, arrival string) *Order {
	return &Order{
		id:      id,
		arrival: aToI(arrival, "order arrival"),
	}
}

type Truck struct {
	id       int
	start    int
	end      int
	capacity int
	orders   int
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

func newTruck(id int, start, end, capacity string) *Truck {
	return &Truck{
		id:       id,
		start:    aToI(start, "truck start"),
		end:      aToI(end, "truck end"),
		capacity: aToI(capacity, "truck capacity"),
	}
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

func aToI(a string, message string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		log.Fatalf("failed to convert string '%s' to int. Error: %+v. Message: %s", a, err, message)
	}
	return i
}

func (ob *OrderBatch) printTrucks() {
	slices.SortFunc(ob.trucks, func(a, b *Truck) int {
		return cmp.Or(
			cmp.Compare(a.id, b.id),
		)
	})
	for _, truck := range ob.trucks {
		fmt.Printf("%d %d %d %d ||", truck.id, truck.start, truck.end, truck.capacity)
	}
	fmt.Println("")
}

func printTrucksStart() {
	for i, truck := range sortedTrucks {
		fmt.Printf("%d:%d ||", i, truck.start)
	}
	fmt.Println("===============")
}

func printOrders(orders []*Order) {
	for i, order := range orders {
		fmt.Printf("order:%d:%d ||", i, order.arrival)
	}
	fmt.Println("===============")
}
