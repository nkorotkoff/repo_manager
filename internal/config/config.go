package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"sync"
)

type Config struct {
	GitEmail            string
	GitPassword         string
	AdminHost           string
	AdminPort           string
	TelegramAccessToken string
	initialized         bool
	mu                  sync.Mutex
}

func (c *Config) LoadConfig() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return nil
	}

	err := godotenv.Load("../.env")
	if err != nil {
		return fmt.Errorf("error loading .env file: %s", err)
	}

	c.GitEmail = os.Getenv("GIT_EMAIL")
	c.GitPassword = os.Getenv("GIT_PASSWORD")
	c.AdminHost = os.Getenv("ADMIN_HOST")
	c.AdminPort = os.Getenv("ADMIN_PORT")
	c.TelegramAccessToken = os.Getenv("TELEGRAM_ACCESS_TOKEN")
	c.initialized = true

	return nil
}

func (c *Config) GetConfig() (Config, error) {
	err := c.LoadConfig()
	if err != nil {
		return Config{}, err
	}
	return *c, nil
}
