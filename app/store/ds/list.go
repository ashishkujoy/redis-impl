package ds

import (
	"errors"
	"fmt"
)

type listNode struct {
	value  string
	parent *listNode
	child  *listNode
}

type List struct {
	length int
	head   *listNode
	tail   *listNode
}

func NewList(value string) *List {
	node := &listNode{value: value}
	return &List{
		head:   node,
		tail:   node,
		length: 1,
	}
}

func (l *List) Len() int {
	return l.length
}

func (l *List) RPush(value string) int {
	node := &listNode{value: value}
	l.tail.child = node
	node.parent = l.tail
	l.tail = node
	l.length++
	return l.Len()
}

func (l *List) LPush(value string) int {
	node := &listNode{value: value}
	l.head.parent = node
	node.child = l.head
	l.head = node
	l.length++
	return l.Len()
}

func (l *List) getNodeByIndex(index int) *listNode {
	head := l.head
	for i := 0; i < index; i++ {
		if head == nil {
			return nil
		}
		head = head.child
	}
	return head
}

func makeStartAndEnd(length int, start int, end int) (int, int) {
	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}
	return max(0, start), min(end, length-1)
}

func (l *List) LRange(start int, end int) []string {
	start, end = makeStartAndEnd(l.length, start, end)
	startingNode := l.getNodeByIndex(start)
	if startingNode == nil {
		return make([]string, 0)
	}
	var elements []string
	for i := start; i <= end; i++ {
		elements = append(elements, startingNode.value)
		startingNode = startingNode.child
	}
	return elements
}

func (l *List) LPop() (string, error) {
	oldHead := l.head
	if oldHead == nil {
		return "", errors.New("list is empty")
	}
	newHead := oldHead.child
	newHead.parent = nil
	l.head = newHead
	l.length--
	return oldHead.value, nil
}
