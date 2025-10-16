package ds

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
