package storage

import (
	"context"
	"sync"
	"time"
)

type ClicksBulkInserter interface {
	BulkInsertClicks(ctx context.Context, clicks map[int]int) error
}

// ClickCounter Структура для хранения батчей кликов in-memory
type ClickCounter struct {
	mu         sync.RWMutex
	counts     map[int]int
	shutdownCh chan struct{}
	storage    ClicksBulkInserter
}

func NewClickCounter(storage ClicksBulkInserter) *ClickCounter {
	cc := &ClickCounter{
		mu:         sync.RWMutex{},
		counts:     make(map[int]int),
		shutdownCh: make(chan struct{}),
		storage:    storage,
	}

	go cc.syncWorker()
	return cc
}

func (cc *ClickCounter) AddClick(bannerID int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.counts[bannerID]++
}

func (cc *ClickCounter) syncWorker() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	sizeCheckTicker := time.NewTicker(30 * time.Second)
	defer sizeCheckTicker.Stop()

	for {
		select {
		case <-ticker.C:
			//todo: хендлинг ошибок
			cc.flushToDB()
		case <-cc.shutdownCh:
			cc.flushToDB()
			return
		case <-sizeCheckTicker.C:
			if len(cc.counts) > 1000 {
				cc.flushToDB()
			}
		}
	}
}

func (cc *ClickCounter) flushToDB() error {
	cc.mu.Lock()
	if len(cc.counts) == 0 {
		cc.mu.Unlock()
		return nil
	}

	//copy the cc.counts map
	counts := make(map[int]int)
	for k, v := range cc.counts {
		counts[k] = v
	}

	cc.counts = make(map[int]int)
	cc.mu.Unlock()

	err := cc.storage.BulkInsertClicks(context.Background(), counts)
	if err != nil {
		return err
	}
	return nil
}

func (cc *ClickCounter) Shutdown() {
	close(cc.shutdownCh)
}
