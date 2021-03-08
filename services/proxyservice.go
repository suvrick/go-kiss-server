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
	var err error
	for _, v := range urls {
		p := *model.NewProxy(v)
		p.ID, err = s.proxyRepository.Create(p)
		if err == nil {
			proxies = append(proxies, p)
		}
	}

	return proxies, err
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

// DeleteAll ...
func (s *ProxyService) DeleteAll() error {
	return s.proxyRepository.DeleteAll()
}
