package repositories

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type Proxy struct {
	gorm.Model

	Host     string `json:"host"`
	Port     string `json:"port"`
	Login    string `gorm:"primaryKey" json:"login"`
	Password string `json:"password"`

	UseToday bool   `json:"use_today"`
	IsError  bool   `json:"is_error"`
	ErrorMsg string `json:"error_msg"`
}

func NewProxy(url string) (*Proxy, error) {
	array := strings.Split(url, ":")

	if len(array) != 4 {
		return nil, errors.New("invalid entres string")
	}

	return &Proxy{
		Host:     array[0],
		Port:     array[1],
		Login:    array[2],
		Password: array[3],
	}, nil
}

type ProxyRepository struct {
	db     *gorm.DB
	locker sync.Mutex
}

// NewBotRepository ...
func NewProxyRepository(db *gorm.DB) *ProxyRepository {
	return &ProxyRepository{
		db:     db,
		locker: sync.Mutex{},
	}
}

func (repo *ProxyRepository) Add(proxy *Proxy) error {
	repo.locker.Lock()
	defer repo.locker.Unlock()

	return repo.db.Table("proxies").Create(proxy).Error
}

func (repo *ProxyRepository) Get() (*Proxy, error) {
	repo.locker.Lock()
	defer repo.locker.Unlock()

	p := &Proxy{}
	err := repo.db.Table("proxies").Where("use_today = ? AND is_error = ?", false, false).First(p).Error
	if err != nil {
		return p, err
	}

	p.UseToday = true
	err = repo.db.Table("proxies").Save(p).Error
	return p, err
}

func (repo *ProxyRepository) GetString() (string, error) {
	p, err := repo.Get()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v:%v:%v:%v", p.Host, p.Port, p.Login, p.Password), nil
}

func (repo *ProxyRepository) UpdateString(login string, isError bool) error {
	repo.locker.Lock()
	defer repo.locker.Unlock()

	p := &Proxy{}
	repo.db.Table("proxies").First(p, "login = ?", login)
	p.IsError = isError

	return repo.Update(p)
}

func (repo *ProxyRepository) Update(proxy *Proxy) error {
	repo.locker.Lock()
	defer repo.locker.Unlock()

	return repo.db.Table("proxies").Save(proxy).Error
}

func (repo *ProxyRepository) Delete(proxy *Proxy) error {
	repo.locker.Lock()
	defer repo.locker.Unlock()

	return nil
}

func (repo *ProxyRepository) ResetProxy() {
	all := make([]Proxy, 0)
	repo.db.Table("proxies").Where("host != ?", "").Scan(&all)

	for _, p := range all {
		p.UseToday = false
		p.IsError = false
		repo.db.Table("proxies").Save(p)
	}
}
