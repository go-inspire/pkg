package config

import (
	"bufio"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	path  string            //配置文件路径
	fw    *fsnotify.Watcher //监控文件变化，自动加载
	value interface{}
}

func NewLoader(path string, v interface{}) *Config {
	path, _ = filepath.Abs(path)
	return &Config{
		path:  path,
		value: v,
	}
}

func (c *Config) Start() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	c.fw = watcher
	err = watcher.Add(c.path)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			log.Print("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Print("modified file:", event.Name)
				err := c.Load()
				if err != nil {
					log.Print("error:", err)
				}

			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return err
			}
			log.Print("error:", err)
		}
	}

}

func (c *Config) Stop() error {
	log.Print("Stop config file watcher...")
	return c.fw.Close()
}

func (c *Config) Load() error {
	log.Printf("Loading config: %v", c.path)
	err := read(c.path, c.value)
	if err != nil {
		panic("Config error: " + err.Error())
	}
	return err
}

func read(file string, v interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = yaml.NewDecoder(bufio.NewReader(f)).Decode(v)
	return err
}
