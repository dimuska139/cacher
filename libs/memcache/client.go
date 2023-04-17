package memcache

import (
	"bufio"
	"bytes"
	"fmt"
	"runtime"
	"sync"
)

// Client клиент для работы с Memcache
type Client struct {
	cfg      *Config
	mx       sync.Mutex
	connPool *Pool
}

// NewMemcacheClient создает MemcacheClient
func NewMemcacheClient(cfg *Config) *Client {
	client := &Client{
		cfg:      cfg,
		connPool: NewPool(cfg),
	}
	runtime.SetFinalizer(client, finalizer)
	return client
}

// finalizer вызывается сборщиком мусора для корректного завершения работы MemcacheClient
func finalizer(c *Client) {
	c.connPool.closeAllConnections()
}

// Get получает запись из Memcache
func (c *Client) Get(key string) ([]byte, error) {
	serverAddress := c.connPool.GetServerAddr(key)
	conn, err := c.connPool.AcquireConnection(serverAddress)
	if err != nil {
		return nil, fmt.Errorf("can't get connection from pool: %w", err)
	}

	defer c.connPool.ReleaseConnection(serverAddress, conn)

	buf := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	if _, err = fmt.Fprintf(buf, "gets %s\r\n", key); err != nil {
		return nil, fmt.Errorf("can't format command and write bytes: %w", err)
	}

	if err := buf.Flush(); err != nil {
		return nil, fmt.Errorf("can't write buffered data to io.Writer: %w", err)
	}

	for {
		row, err := buf.ReadSlice('\n')
		if err != nil {
			return nil, fmt.Errorf("can't read slice: %w", err)
		}

		if bytes.HasPrefix(row, []byte("VALUE ")) {
			continue
		}

		if string(row) == "END\r\n" {
			break
		}

		if !bytes.HasSuffix(row, []byte("\r\n")) {
			return nil, fmt.Errorf("invalid data in cache: %s", string(row))
		}

		return row[0 : len(row)-2], nil // Удаляем \r\n в конце строки
	}
	return nil, nil
}

// Set делает запись в Memcache
func (c *Client) Set(key string, value []byte, expiration int64) error {
	serverAddress := c.connPool.GetServerAddr(key)

	conn, err := c.connPool.AcquireConnection(serverAddress)
	if err != nil {
		return fmt.Errorf("can't get connection from pool: %w", err)
	}

	defer c.connPool.ReleaseConnection(serverAddress, conn)

	buf := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	if _, err = fmt.Fprintf(buf, "set %s 0 %d %d\r\n", key, expiration, len(value)); err != nil {
		return fmt.Errorf("can't format command and write bytes: %w", err)
	}

	if _, err := buf.Write(append(value, []byte("\r\n")...)); err != nil {
		return fmt.Errorf("can't write bytes: %w", err)
	}

	if err := buf.Flush(); err != nil {
		return fmt.Errorf("can't write buffered data to io.Writer: %w", err)
	}

	row, err := buf.ReadSlice('\n')
	if err != nil {
		return fmt.Errorf("can't read slice: %w", err)
	}
	if string(row) == "STORED\r\n" {
		return nil
	}
	return fmt.Errorf("can't store data: %s", string(row))
}

// Delete удаляет запись из Memcache
func (c *Client) Delete(key string) error {
	serverAddress := c.connPool.GetServerAddr(key)
	conn, err := c.connPool.AcquireConnection(serverAddress)
	if err != nil {
		return fmt.Errorf("can't get connection from pool: %w", err)
	}

	defer c.connPool.ReleaseConnection(serverAddress, conn)

	buf := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	if _, err = fmt.Fprintf(buf, "delete %s\r\n", key); err != nil {
		return fmt.Errorf("can't format command and write bytes: %w", err)
	}

	if err := buf.Flush(); err != nil {
		return fmt.Errorf("can't write buffered data to io.Writer: %w", err)
	}

	row, err := buf.ReadSlice('\n')
	if err != nil {
		return fmt.Errorf("can't read slice: %w", err)
	}

	if string(row) == "DELETED\r\n" || string(row) == "NOT_FOUND\r\n" {
		return nil
	}

	return fmt.Errorf("can't delete item: %s", string(row))
}
