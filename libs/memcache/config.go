package memcache

import (
	"net"
	"time"
)

const (
	DefaultTimeout  = time.Second * 200
	DefaultPoolSize = 5
)

// Config конфигурация для библиотеки-клиента Memcache
type Config struct {
	timeout  time.Duration
	poolSize int
	servers  []net.Addr
}

// NewConfig создаёт конфигурацию для библиотеки-клиента Memcache
func NewConfig(servers []net.Addr, poolSize int, timeout time.Duration) *Config {
	return &Config{
		timeout:  timeout,
		poolSize: poolSize,
		servers:  servers,
	}
}

// Timeout возвращает таймаут
func (c *Config) Timeout() time.Duration {
	if c.timeout >= 0 {
		return c.timeout
	}
	return DefaultTimeout
}

// PoolSize возвращает размер пула соединений
func (c *Config) PoolSize() int {
	if c.poolSize >= 0 {
		return c.poolSize
	}
	return DefaultPoolSize
}
