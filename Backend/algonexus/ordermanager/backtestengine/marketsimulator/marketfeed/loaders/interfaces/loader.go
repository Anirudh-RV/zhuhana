package interfaces

import "algonexus/ordermanager/backtestengine/models"

type TickIterator interface {
	Next() (*models.MarketTick, error)
	Close() error
}

type TickLoader struct {
	iterator *TickIterator
}
