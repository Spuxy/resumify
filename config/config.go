package config

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	Dir    string
	Values map[string]string
}

func New(path string) (*Config, error) {
	var cfg Config
	values := make(map[string]string)
	cfg.Dir = path
	cfg.Values = values

	err := (&cfg).readAll()
	if err != nil {
		return &Config{}, err
	}

	return &cfg, nil
}

func (c *Config) readAll() error {

	files, err := ioutil.ReadDir(c.Dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := c.read(file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) read(file string) error {
	content, err := os.Open(fmt.Sprintf("%s/%s", c.Dir, file))
	defer content.Close()

	if err != nil {
		return err
	}

	buffer := bufio.NewScanner(content)

	for buffer.Scan() {
		c.set(buffer.Bytes())
	}

	return nil
}

func (c *Config) set(line []byte) {
	key_value := strings.Split(string(line), "=")
	c.Values[key_value[0]] = strings.Trim(key_value[1], `"`)
}

func (c *Config) Get(key string) string {
	return c.Values[key]
}
