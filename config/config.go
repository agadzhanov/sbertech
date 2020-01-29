package config

import "fmt"

type config struct {
	Network string
	Domain  string
	Port    uint64
}

func (c *config) Address() string {
	return fmt.Sprintf("%s:%d", c.Domain, c.Port)
}

func GetConfig() *config {
	return &config{
		Network: "tcp",
		Domain:  "localhost",
		Port:    2250,
	}
}
