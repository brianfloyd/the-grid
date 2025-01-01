package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

var InvalidKey = errors.New("Invalid configuration key specified")

func getEnv(key string) (string, error) {
	val, set := os.LookupEnv(key)
	if !set {
		return "", InvalidKey
	}
	return val, nil
}

type Config struct {
	path   string
	domain string
	values interface{}
}

func NewConfig() *Config {
	domain, err := getEnv("DOMAIN")

	if err != nil {
		domain = "test"
	}

	return &Config{
		path:   fmt.Sprintf("./cfg/%s.json", strings.ToLower(domain)),
		domain: domain,
	}
}

func (c *Config) Read() error {
	bytes, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}

	var data interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	c.values = data
	return nil
}

func (c *Config) ValueOrPanic(key string) string {
	v, err := c.Value(key)
	if err != nil {
		panic("Could not find value for key: " + key)
	}
	return v
}

func (c *Config) Value(key string) (string, error) {
	obj, ok := c.values.(map[string]interface{})
	if !ok {
		return "", InvalidKey
	}

	parts := strings.Split(key, ".")
	for _, part := range parts[:len(parts)-1] {
		i, ok := obj[part]
		if !ok {
			return "", InvalidKey
		}
		obj, ok = i.(map[string]interface{})
		if !ok {
			return "", InvalidKey
		}
	}

	value, ok := obj[parts[len(parts)-1]]
	if !ok {
		return "", InvalidKey
	}

	valueStr, ok := value.(string)
	if !ok {
		return "", InvalidKey
	}

	n := len(valueStr)
	if valueStr[0:2] == "{{" && valueStr[n-2:n] == "}}" {
		var err error
		valueStr, err = getEnv(valueStr[2 : n-2])
		if err != nil {
			return "", InvalidKey
		}
	}

	return valueStr, nil
}
