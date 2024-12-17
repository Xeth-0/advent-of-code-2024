package main

import (
	. "advent2024/utils/vector"
)

/*
Man I'm such a noob. Don't even know how to use a heap. Welp, here goes.

Go's heap package ("container/heap") requires any data structure that is to be used as a heap
to implement the heap.Interface ...interface. Doing this in go doesn't entail inheritance or any
of that shit, it just means the data structure(the struct) should have some methods/properties.

In this case, to implement heap.Interface, the priority queue needs 5 methods.
	Len()
	Less(i, j int) bool
	Swap(i, j int)
	Push(x, interface{})
	Pop() interface{}

(interface{} is just the placeholder for the generic. Since go doesn't(or at least didn't) have generics)
*/

// For the priority queue (UCS)
type Node struct {
	pos    Vector // Position in cartesian space(for the grid)
	cost   int    // Path Cost
	index  int    // Index in the heap
	dir    string // The direction the robot is facing
	parent *Node  // To find the path taken at the end
}

type PriorityQueue []*Node

func (queue *PriorityQueue) Len() int {
	return len(*queue)
}

func (queue *PriorityQueue) Less(i int, j int) bool {
	return (*queue)[i].cost < (*queue)[j].cost
}

func (queue *PriorityQueue) Swap(i int, j int) {
	(*queue)[i], (*queue)[j] = (*queue)[j], (*queue)[i]
	(*queue)[i].index, (*queue)[j].index = i, j
}

func (queue *PriorityQueue) Push(x interface{}) {
	node := x.(*Node) // this is a type assertion. we're just casting the x interface{}(generic) to a *Node pointer
	node.index = len(*queue)
	*queue = append(*queue, node)
}

func (queue *PriorityQueue) Pop() interface{} {
	old := *queue
	n := len(old)

	// Get then Invalidate the last element
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*queue = old[:n-1]
	return node
}
