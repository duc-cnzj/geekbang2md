package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Cache struct{}

var getwd, _ = os.Getwd()
var cacheDirName = "geekbang-cache-data"
var cacheDir = filepath.Join(getwd, cacheDirName)

func init() {
	os.MkdirAll(cacheDir, 0755)
}

func Dir() string {
	return cacheDir
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
	return fmt.Sprintf(filepath.Join(cacheDirName, "cache-%v.json"), key)
}
