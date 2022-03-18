package image

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/DuC-cnZj/geekbang2md/constant"
	"github.com/DuC-cnZj/geekbang2md/waiter"
)

type Manager struct {
	sync.RWMutex
	images  map[string]string
	baseDir string
	waiter  waiter.Interface
}

func NewManager(baseDir string) *Manager {
	os.MkdirAll(filepath.Join(baseDir, "mp3"), 0755)

	return &Manager{
		RWMutex: sync.RWMutex{},
		images:  map[string]string{},
		baseDir: baseDir,
		waiter:  waiter.NewSigWaiter(constant.ImageDownloadParallelNum),
	}
}

func (m *Manager) Download(u string) (string, error) {
	if path := m.Get(u); path != "" {
		return path, nil
	}
	parse, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("%w path: %s", err, u)
	}
	split := strings.Split(parse.Path, "/")
	name := split[len(split)-1]
	var p string
	if strings.HasSuffix(name, ".mp3") {
		p = filepath.Join(m.baseDir, "mp3", name)
	} else {
		p = filepath.Join(m.baseDir, name)
	}
	stat, err := os.Stat(p)
	if err == nil && stat.Mode().IsRegular() {
		m.Add(u, p)
		return p, nil
	}
	m.waiter.Wait(context.TODO())
	defer m.waiter.Release()
	c := &http.Client{}
	res, err := c.Get(u)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	all, _ := io.ReadAll(res.Body)
	if err := os.WriteFile(p, all, 0644); err != nil {
		return "", fmt.Errorf("err: %w, origin path: %s, write path: %s", err, u, p)
	}
	m.Add(u, p)
	return p, nil
}

func (m *Manager) Has(url string) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.images[url]
	return ok
}

func (m *Manager) Get(url string) string {
	m.RLock()
	defer m.RUnlock()
	path, ok := m.images[url]
	if !ok {
		return ""
	}
	return path
}

func (m *Manager) Add(url, path string) {
	m.Lock()
	defer m.Unlock()
	m.images[url] = path
}
