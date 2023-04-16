package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	// Порт, на котором запустится GRPC-сервер
	GrpcPort int `yaml:"grpc_port"`
	// Уровни логирования (debug, info, warn, error)
	Loglevel string `yaml:"loglevel"`
	// Тип используемого хранилища (memcache или любое другое значения для использования встроенного кеша)
	Storage string `yaml:"storage"`
	// Список серверов Memcache (при использовании storage != memcache можно не указывать)
	// Для упрощения тут поддерживается только TCP, unix-сокеты - нет
	MemcacheServers []string `yaml:"memcache_servers"`
}

// NewConfig инициализирует конфиг
func NewConfig(configPath string) (*Config, error) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("can't read config file: %w", err)
	}

	var cfg Config

	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, fmt.Errorf("can't parse config: %w", err)
	}

	return &cfg, nil
}
