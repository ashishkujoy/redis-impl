package ds

import "fmt"

type Lists struct {
	lists map[string]*List
}

func NewLists() *Lists {
	return &Lists{lists: make(map[string]*List)}
}

func (l *Lists) RPush(name string, value string) int {
	list, ok := l.lists[name]
	if !ok {
		fmt.Printf("List %s is nil\n", name)
		list = NewList(value)
		l.lists[name] = list
		return 1
	}
	fmt.Printf("List %s is exist\n", name)
	return list.RPush(value)
}
