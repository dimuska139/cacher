package memcache

import (
	"fmt"
	"hash/crc32"
	"net"
	"sync"
	"time"
)

// Pool отвечает за пулл соединений с Memcache
type Pool struct {
	mx                   sync.Mutex
	availableConnections map[string][]net.Conn
	cfg                  *Config
}

// NewPool создаёт пулл соедений с Memcache
func NewPool(cfg *Config) *Pool {
	return &Pool{
		availableConnections: make(map[string][]net.Conn),
		cfg:                  cfg,
	}
}

// ReleaseConnection возвращает соединение в пулл
func (c *Pool) ReleaseConnection(addr net.Addr, conn net.Conn) {
	c.mx.Lock()
	defer c.mx.Unlock()

	availableServerConn := c.availableConnections[addr.String()]
	// Если размер предельный размер пула достигнут, то просто закрываем соединение, иначе возвращаем в пул
	if len(availableServerConn) < c.cfg.PoolSize() {
		c.availableConnections[addr.String()] = append(availableServerConn, conn)
	} else {
		conn.Close()
	}
}

// releaseAllConnections закрывает все соединения
func (c *Pool) closeAllConnections() {
	for _, addr := range c.cfg.servers {
		for _, conn := range c.availableConnections[addr.String()] {
			conn.Close()
		}
	}
	c.availableConnections = make(map[string][]net.Conn)
}

// connect устанавливает новое соединение
func (c *Pool) connect(addr net.Addr) (net.Conn, error) {
	conn, err := net.DialTimeout(addr.Network(), addr.String(), c.cfg.Timeout())
	if err == nil {
		return conn, nil
	}

	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		return nil, fmt.Errorf("connection timeout deadline exceeded (%s): %w", addr.String(), err)
	}

	return nil, fmt.Errorf("can't connect to the address %s: %w", addr.String(), err)
}

// getFreeConnection возвращает доступное соединение, если оно есть
func (c *Pool) getFreeConnection(addr net.Addr) (net.Conn, error) {
	availableServerConn, ok := c.availableConnections[addr.String()]
	if !ok || len(availableServerConn) == 0 {
		return nil, nil
	}

	conn := availableServerConn[0]
	c.availableConnections[addr.String()] = availableServerConn[1:]
	err := conn.SetDeadline(time.Now().Add(c.cfg.Timeout()))
	if err != nil {
		return nil, fmt.Errorf("can't set new connection deadline: %w", err)
	}

	return conn, nil
}

// GetServerAddr возвращает адрес сервера Memcache
func (c *Pool) GetServerAddr(key string) net.Addr {
	// См. https://habr.com/ru/articles/42972/
	return c.cfg.servers[crc32.Checksum([]byte(key), crc32.MakeTable(crc32.IEEE))%uint32(len(c.cfg.servers))]
}

// AcquireConnection получает соединение из пула
func (c *Pool) AcquireConnection(addr net.Addr) (net.Conn, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	// Ищем доступное соединение
	existingConnection, err := c.getFreeConnection(addr)
	if err != nil {
		return nil, fmt.Errorf("can't get free connection: %w", err)
	}

	if existingConnection != nil {
		return existingConnection, nil
	}

	// Если не нашли доступное, то открываем новое
	newConn, err := c.connect(addr)
	if err != nil {
		return nil, fmt.Errorf("can't create new connection: %w", err)
	}

	return newConn, nil
}
