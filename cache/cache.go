package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Cache struct{}

var dir string

func Init(baseDir string) {
	dir = filepath.Join(baseDir, ".cache")
	os.MkdirAll(dir, 0755)
}

func Dir() string {
	return dir
}

func (c *Cache) Get(key string) ([]byte, error) {
	file, err := os.ReadFile(c.cachePath(key))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (c *Cache) Set(key string, data interface{}) error {
	marshal, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	if len(marshal) > 0 {
		if err := os.WriteFile(c.cachePath(key), marshal, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache) cachePath(key string) string {
	return fmt.Sprintf(filepath.Join(dir, "cache-%v.json"), key)
}
