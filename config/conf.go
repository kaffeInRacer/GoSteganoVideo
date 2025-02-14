package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	serverConfig
}

type serverConfig struct {
	Host string
	Port string
	Base string
}

func (c *Config) LoadEnvironment() error {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	c.Host = os.Getenv("HOST")
	c.Port = os.Getenv("PORT")
	c.Base = fmt.Sprintf("%s:%s", c.Host, c.Port)
	return nil
}
