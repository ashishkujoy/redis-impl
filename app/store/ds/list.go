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

func (l *List) LRange(start int, end int) []string {
	start = max(start, 0)
	end = min(end, l.Len()-1)
	startingNode := l.getNodeByIndex(start)
	if startingNode == nil {
		return make([]string, 0)
	}
	elements := make([]string, 0, end-start)
	for i := start; i <= end; i++ {
		elements = append(elements, startingNode.value)
		startingNode = startingNode.child
	}
	return elements
}
