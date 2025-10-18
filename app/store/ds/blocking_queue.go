package ds

import "sync"

type BlockedClient struct {
	Key      string
	WakeChan chan string
}

type BlockingQueueManager struct {
	mutex  sync.RWMutex
	queues map[string][]*BlockedClient
}

func NewBlockingQueueManager() *BlockingQueueManager {
	return &BlockingQueueManager{
		mutex:  sync.RWMutex{},
		queues: make(map[string][]*BlockedClient),
	}
}

func (bm *BlockingQueueManager) BlockOn(key string) *BlockedClient {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()
	client := &BlockedClient{Key: key, WakeChan: make(chan string, 1)}
	bm.queues[key] = append(bm.queues[key], client)
	return client
}

func (bm *BlockingQueueManager) AnyBlockOn(key string) bool {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()
	clients, ok := bm.queues[key]
	if !ok {
		return false
	}
	return len(clients) > 0
}

func (bm *BlockingQueueManager) Unblock(key string, value string) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()
	blockedClients, ok := bm.queues[key]
	if !ok {
		return
	}
	firstClient := blockedClients[0]
	remainingClients := blockedClients[1:]

	if len(remainingClients) == 0 {
		delete(bm.queues, key)
	} else {
		bm.queues[key] = remainingClients
	}

	firstClient.WakeChan <- value
}
