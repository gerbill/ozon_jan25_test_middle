package main

import (
	"bytes"
	"cmp"
	"slices"
	"strconv"
)

func (ob *OrderBatch) report() string {
	var report bytes.Buffer
	sortedTrucks = ob.sortedTrucks()
	sortedOrders := ob.sortedOrders()
	for _, order := range sortedOrders {
		order.findTruck()
	}
	// orders := sortedOrders
	// for _, truck := range sortedTrucks {
	// 	if len(orders) == 0 {
	// 		break
	// 	}
	// 	orders = truck.findOrders(orders)
	// }
	for _, order := range ob.orders2 {
		report.WriteString(strconv.Itoa(order.truckId) + " ")
	}
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
