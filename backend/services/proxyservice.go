package services

import (
	"sync"

	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/repositories"
)

// ProxyService ...
type ProxyService struct {
	proxyRepository *repositories.ProxyRepository
	locker          sync.Mutex
}

// NewProxyService ...
func NewProxyService(repo *repositories.ProxyRepository) *ProxyService {
	return &ProxyService{
		proxyRepository: repo,
	}
}

// AddRange ...
func (s *ProxyService) AddRange(urls []string) ([]model.Proxy, error) {

	proxies := make([]model.Proxy, 0)

	for _, v := range urls {
		p := *model.NewProxy(v)
		p.ID, _ = s.proxyRepository.Create(p)
		if p.ID != 0 {
			proxies = append(proxies, p)
		}
	}

	return proxies, nil
}

// Free ...
func (s *ProxyService) Free() (model.Proxy, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	return s.proxyRepository.Free()
}

// Update ...
func (s *ProxyService) Update(proxy model.Proxy) error {
	s.locker.Lock()
	defer s.locker.Unlock()

	return s.proxyRepository.Update(proxy)
}

// All ...
func (s *ProxyService) All() ([]model.Proxy, error) {
	return s.proxyRepository.All()
}

// Clear ...
func (s *ProxyService) Clear() error {
	return s.proxyRepository.Clear()
}
