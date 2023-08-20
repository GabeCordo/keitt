package processor

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func NewConfig(name string) *Config {
	config := new(Config)
	config.Name = name
	config.Net.Host = "localhost"
	config.Net.Port = 5023
	config.Debug = true
	return config
}

func ConfigFromYAML(config *Config, path string) error {
	if _, err := os.Stat(path); err != nil {
		// file does not exist
		log.Println(err)
		return err
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		// error reading the file
		log.Println(err)
		return err
	}

	err = yaml.Unmarshal([]byte(file), config)
	if err != nil {
		// the file is not a JSON or is a malformed (fields missing) core
		log.Println(err)
		return err
	}

	return nil
}

var (
	configLock     = &sync.Mutex{}
	configInstance *Config
)
