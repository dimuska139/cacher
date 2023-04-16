package memcache

import (
	"fmt"
	memcacheClient "github.com/dimuska139/cacher/libs/memcache"
	"github.com/dimuska139/cacher/pkg/config"
	"net"
	"time"
)

// NewClient инициирует библиотеку для работы с Memcache
func NewClient(config *config.Config) (*memcacheClient.Client, error) {
	srvs := make([]net.Addr, 0, len(config.MemcacheServers))
	// Для упрощения тут поддерживается только TCP, unix-сокеты - нет
	for _, addr := range config.MemcacheServers {
		tcpaddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			return nil, fmt.Errorf("can't resolve TCP address %s: %w", addr, err)
		}
		srvs = append(srvs, tcpaddr)
	}

	memcacheClientConfig := memcacheClient.NewConfig(srvs, 5, time.Second)
	return memcacheClient.NewMemcacheClient(memcacheClientConfig), nil
}
