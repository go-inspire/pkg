package config

import (
	"bufio"
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/go-inspire/pkg/log"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Loader struct {
	path  string            //配置文件路径
	fw    *fsnotify.Watcher //监控文件变化，自动加载
	value interface{}
}

func (c *Loader) Value() interface{} {
	return c.value
}

func NewLoader(path string, v interface{}) *Loader {
	path, _ = filepath.Abs(path)
	return &Loader{
		path:  path,
		value: v,
	}
}

func (c *Loader) Start(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	c.fw = watcher
	err = watcher.Add(c.path)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case event := <-watcher.Events:
			log.Info("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op == fsnotify.Rename {
				log.Info("modified file:", event.Name)
				err := c.Load()
				if err != nil {
					log.Error("error:", err)
				}

			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return err
			}
			log.Error("error:", err)
		}
	}

}

func (c *Loader) Stop(ctx context.Context) error {
	log.Info("Stop config file watcher...")
	return c.fw.Close()
}

func (c *Loader) Load() error {
	log.Printf("Loading config: %v", c.path)
	err := read(c.path, c.value)
	if err != nil {
		panic("Loader error: " + err.Error())
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
