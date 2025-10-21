package ds

import (
	"context"
	"sync"
)

type BlockedClient struct {
	ClientId    int
	Key         string
	WakeChan    chan string
	TimeoutChan chan string
}

func NewBlockedClient(clientId int, key string) *BlockedClient {
	wakeChan := make(chan string)
	timeoutChan := make(chan string)

	return &BlockedClient{
		ClientId:    clientId,
		Key:         key,
		WakeChan:    wakeChan,
		TimeoutChan: timeoutChan,
	}
}

type BlockingQueueManager struct {
	mutex   sync.RWMutex
	idCount int
	queues  map[string][]*BlockedClient
}

func NewBlockingQueueManager() *BlockingQueueManager {
	return &BlockingQueueManager{
		mutex:  sync.RWMutex{},
		queues: make(map[string][]*BlockedClient),
	}
}

func (bm *BlockingQueueManager) generateClientId() int {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()
	id := bm.idCount
	bm.idCount++
	return id
}

func (bm *BlockingQueueManager) BlockOn(key string, ctx context.Context) *BlockedClient {
	clientId := bm.generateClientId()
	bm.mutex.Lock()
	defer bm.mutex.Unlock()
	client := NewBlockedClient(clientId, key)
	go bm.removeOnCancel(ctx, clientId, key)
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

func (bm *BlockingQueueManager) removeOnCancel(ctx context.Context, id int, key string) {
	_ = <-ctx.Done()
	bm.mutex.Lock()
	defer bm.mutex.Unlock()
	queues := bm.queues[key]

	for _, client := range queues {
		if client.ClientId == id {
			client.TimeoutChan <- "timed out"
		} else {
			queues = append(queues, client)
		}
	}
	bm.queues[key] = queues
}
