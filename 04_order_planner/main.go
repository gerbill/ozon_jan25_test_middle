package main

import (
	"bufio"
	"cmp"
	"fmt"
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
		maxWindowDuration = 0
		lastAvailableTruckId = 0
		// ob.printTrucks()
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

func printTrucks(trucks []*Truck) {
	slices.SortFunc(trucks, func(a, b *Truck) int {
		return cmp.Or(
			cmp.Compare(a.id, b.id),
		)
	})
	for i, truck := range trucks {
		fmt.Printf("%+v ||", truck)
		if i > 10 {
			break
		}
	}
	fmt.Println("")
}

func printTrucksStart() {
	for i, truck := range sortedTrucks {
		fmt.Printf("%d:%d ||", i, truck.start)
		if i > 10 {
			break
		}
	}
	fmt.Println("================")
}

func printOrders(orders []*Order) {
	for i, order := range orders {
		fmt.Printf("%+v ||", order)
		if i > 10 {
			break
		}
	}
	fmt.Println("================")
}
