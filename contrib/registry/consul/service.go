package consul

import (
	"context"
	"sync"
	"sync/atomic"

	"gitlab.wwgame.com/wwgame/kratos/v2/registry"
)

type serviceSet struct {
	serviceName string
	watcher     map[*watcher]struct{}
	services    *atomic.Value
	lock        sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
}

func (s *serviceSet) broadcast(ss []*registry.ServiceInstance) {
	s.services.Store(ss)
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k := range s.watcher {
		select {
		case k.event <- struct{}{}:
		default:
		}
	}
}
