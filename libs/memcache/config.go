package memcache

import (
	"net"
	"time"
)

const (
	DefaultTimeout  = time.Second * 200
	DefaultPoolSize = 5
)

type Config struct {
	timeout  time.Duration
	poolSize int
	servers  []net.Addr
}

func NewConfig(servers []net.Addr, poolSize int, timeout time.Duration) *Config {
	return &Config{
		timeout:  timeout,
		poolSize: poolSize,
		servers:  servers,
	}
}

func (c *Config) Timeout() time.Duration {
	if c.timeout >= 0 {
		return c.timeout
	}
	return DefaultTimeout
}

func (c *Config) PoolSize() int {
	if c.poolSize >= 0 {
		return c.poolSize
	}
	return DefaultPoolSize
}
