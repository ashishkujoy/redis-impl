package ds

import (
	"fmt"
	"sync"
)

type Lists struct {
	mutex                sync.RWMutex
	lists                map[string]*List
	blockingQueueManager *BlockingQueueManager
}

func NewLists(blockingQueueManager *BlockingQueueManager) *Lists {
	return &Lists{
		mutex:                sync.RWMutex{},
		lists:                make(map[string]*List),
		blockingQueueManager: blockingQueueManager,
	}
}

func (l *Lists) RPush(name string, values []string) int {
	list, ok := l.lists[name]
	if !ok {
		list = NewList(values[0])
		l.lists[name] = list
		values = values[1:]
	}
	for _, value := range values {
		list.RPush(value)
	}
	go l.Wake(name)
	return list.length
}

func (l *Lists) Wake(key string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.blockingQueueManager.AnyBlockOn(key) {
		values, _ := l.LPop(key, 1)
		l.blockingQueueManager.Unblock(key, values[0])
	}
	return
}

func (l *Lists) LPush(name string, values []string) int {
	list, ok := l.lists[name]
	if !ok {
		fmt.Printf("List %s is nil\n", name)
		list = NewList(values[0])

		l.lists[name] = list
		values = values[1:]
	}
	fmt.Printf("List %s is exist\n", name)
	for _, value := range values {
		list.LPush(value)
	}
	go l.Wake(name)
	return list.length
}

func (l *Lists) LRange(key string, start int, end int) []string {
	list, ok := l.lists[key]
	if !ok {
		return make([]string, 0)
	}
	return list.LRange(start, end)
}

func (l *Lists) LLen(key string) int {
	list, ok := l.lists[key]
	if !ok {
		return 0
	}
	return list.length
}

func (l *Lists) LPop(key string, count int) ([]string, error) {
	list, ok := l.lists[key]
	if !ok {
		return nil, fmt.Errorf("List %s is not present\n", key)
	}
	return list.LPop(count)
}
