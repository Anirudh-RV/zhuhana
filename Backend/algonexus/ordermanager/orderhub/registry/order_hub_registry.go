package registry

import (
	"algonexus/logger"
	"algonexus/ordermanager/orderhub/runtime"
	"sync"
)

type OrderHubRegistry struct {
	logger *logger.Logger
	orders map[string]*runtime.OrderSession
	rwmu   sync.RWMutex
}

func NewOrderHubRegistry(logger *logger.Logger) *OrderHubRegistry {
	return &OrderHubRegistry{
		logger: logger,
		orders: make(map[string]*runtime.OrderSession),
	}
}

func (r *OrderHubRegistry) Get(id string) *runtime.OrderSession {
	r.rwmu.RLock()
	defer r.rwmu.RUnlock()
	handle, ok := r.orders[id]
	if !ok {
		return nil
	}
	return handle
}

func (r *OrderHubRegistry) Update(id string, handle *runtime.OrderSession) {
	r.rwmu.Lock()
	defer r.rwmu.Unlock()
	r.orders[id] = handle
}

func (r *OrderHubRegistry) Delete(id string) {
	r.rwmu.Lock()
	defer r.rwmu.Unlock()
	delete(r.orders, id)
}
