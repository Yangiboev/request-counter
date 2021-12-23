package server

import (
	"fmt"
	"time"
)

type RequestCount struct {
	Timestamp   time.Time
	Count       int
	Accumulated int
}

func (r RequestCount) Empty() bool {
	return r.Timestamp.IsZero()
}

func (r RequestCount) Dump() string {
	return fmt.Sprintf("{timestamp:%v, requestsCount:%v, accumulatedRequestCount:%v}", r.Timestamp.String(), r.Count, r.Accumulated)
}

func (r *RequestCount) Increment() {
	r.Count++
	r.Accumulated++
}
func (r RequestCount) CompareTimestampWithPrecision(t time.Time, precision time.Duration) bool {
	return r.Timestamp.Truncate(precision) == t.Truncate(precision)
}

type requestCountNode struct {
	data  RequestCount
	left  *requestCountNode
	right *requestCountNode
}

func (node requestCountNode) WithinDurationBefore(duration time.Duration, precision time.Duration, reference RequestCount) (bool, time.Duration) {
	difference := reference.Timestamp.Sub(node.data.Timestamp)
	return difference.Truncate(precision).Nanoseconds() <= duration.Truncate(precision).Nanoseconds(), difference
}

type RequestCounter struct {
	head *requestCountNode
	tail *requestCountNode
}

type requestCountList []RequestCount

func (r RequestCounter) getNodes() requestCountList {
	nodes := make(requestCountList, 0, 100)
	currentNode := r.head
	for currentNode != nil {
		nodes = append(nodes, currentNode.data)
		currentNode = currentNode.right
	}

	return nodes
}

func (values requestCountList) ToRequestCounter() RequestCounter {
	var list RequestCounter
	for _, value := range values {
		list = list.AppendToTail(value)
	}

	return list
}

func (list RequestCounter) AppendToTail(data RequestCount) RequestCounter {
	newNode := requestCountNode{data: data}
	if list.head == nil {
		list.head = &newNode
		list.tail = &newNode
	} else {
		list.tail.right = &newNode
		newNode.left = list.tail
		list.tail = &newNode
	}

	return list
}

func (list RequestCounter) frontDiscardUntil(lastNodeToDiscard *requestCountNode) RequestCounter {
	currentNode := list.head
	if lastNodeToDiscard == list.tail {
		list.head = nil
		list.tail = nil
	} else {
		list.head = lastNodeToDiscard.right
		list.head.left = nil
	}

	for currentNode != nil {
		atLastNode := false
		if currentNode == lastNodeToDiscard {
			atLastNode = true
		}

		temp := currentNode.right
		currentNode.left = nil
		currentNode.right = nil
		currentNode = temp

		if atLastNode {
			break
		}
	}

	return list
}

func (list RequestCounter) UpdateTotals(reference RequestCount, timeFrame time.Duration, precision time.Duration) RequestCounter {
	currentNode := list.tail
	for currentNode != nil {
		if withinTimeFrame, _ := currentNode.WithinDurationBefore(timeFrame, precision, reference); withinTimeFrame {
			if currentNode.right != nil {
				currentNode.data.Accumulated = currentNode.data.Count + currentNode.right.data.Accumulated
			} else {
				currentNode.data.Accumulated = currentNode.data.Count
			}
			currentNode = currentNode.left
		} else {
			list = list.frontDiscardUntil(currentNode)
			break

		}
	}

	return list
}

func (list RequestCounter) TotalAccumulatedRequestCount() int {
	if list.head != nil {
		return list.head.data.Accumulated
	} else {
		return 0
	}
}
