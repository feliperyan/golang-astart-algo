// Basically a copy-paste of:
// https://golang.org/src/container/heap/example_pq_test.go

package main

import (
	"container/heap"
	"fmt"
)

type MapElement struct {
	pos_x    int
	pos_y    int
	name     string
	passable bool
}

type Map2d struct {
	x     int
	y     int
	two_d [][]MapElement
}

// An Item is something we manage in a priority queue.
type Item struct {
	value    *MapElement // The value of the item; arbitrary.
	priority int         // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// ATTENTION: For the A* Algorithm I want the item with the LOWEST priority (cost)
	// so the sign is inverted from the original code and regular priority queues.
	return pq[i].priority <= pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value *MapElement, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func Example_priorityQueue() {

	// Some items and their priorities.
	i1 := MapElement{2, 2, "a", false}
	i2 := MapElement{2, 3, "b", false}
	i3 := MapElement{2, 4, "c", false}

	items := map[*MapElement]int{
		&i1: 3, &i2: 2, &i3: 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    &MapElement{2, 1, "d", false},
		priority: 1,
	}
	heap.Push(&pq, item)
	// pq.Push(item)
	pq.update(item, item.value, 5)

	// Take the items out; they arrive in INCREASING priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value.name)
	}
	fmt.Println("\n")

}