package zhuanlan

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/DuC-cnZj/geekbang2md/image"
	"github.com/DuC-cnZj/geekbang2md/utils"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/dustin/go-humanize"
)

var imgRegexp = regexp.MustCompile(`!\[(.*?)]\((.*?)\)`)

type MDWriter struct {
	title        string
	baseDir      string
	imageManager *image.Manager
}

func NewMDWriter(baseDir, title string, imgM *image.Manager) *MDWriter {
	os.MkdirAll(baseDir, 0755)
	return &MDWriter{baseDir: baseDir, imageManager: imgM, title: title}
}

func (w *MDWriter) GetFileName(filename string) string {
	filename = utils.FilterCharacters(filename)
	name := filepath.Join(w.baseDir, filename)
	if strings.HasSuffix(name, ".md") {
		return name
	}
	return name + ".md"
}

func (w *MDWriter) FileExists(filename string) (os.FileInfo, bool) {
	st, err := os.Stat(w.GetFileName(filename))
	if err == nil && st.Size() > 0 {
		return st, true
	}
	if os.IsNotExist(err) {
		return nil, false
	}
	return nil, false
}

func (w *MDWriter) WriteFile(audioDownloadURL, audioDubber, audioSize, audioTime, title string, html string) error {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(html)
	if err != nil {
		return err
	}
	var ss = &SafeString{s: markdown}
	//拿出图片，抓图片
	images := FindAllImages(markdown)
	if audioDownloadURL != "" {
		images = append(images, audioDownloadURL)
	}

	wg := sync.WaitGroup{}
	for _, s := range images {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			if s == "" {
				return
			}
			download, err := w.imageManager.Download(s)
			if err != nil {
				log.Println(err)
			} else {
				rel, _ := filepath.Rel(w.baseDir, download)
				ss.Replace(s, rel)
			}
		}(s)
	}
	wg.Wait()

	rel, _ := filepath.Rel(w.baseDir, w.imageManager.Get(audioDownloadURL))
	mdheader := fmt.Sprintf(`
# %s

`, title)
	mdAudio := fmt.Sprintf(`
<span style="font-size: 12px">讲述：%s </span>&nbsp;&nbsp;<span style="font-size: 12px">大小：%s </span>&nbsp;&nbsp;<span style="font-size: 12px">时长：%s</span>

<audio id="audio" controls="" preload="none">
  <source id="mp3" src="%s">
</audio>

`, audioDubber, audioTime, audioSize, rel)
	if w.imageManager.Get(audioDownloadURL) == "" {
		mdAudio = ""
	}
	ss.Set(mdheader + mdAudio + ss.Get())
	file, err := os.OpenFile(w.GetFileName(title), os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write([]byte(ss.Get())); err != nil {
		return err
	}
	log.Printf("[WRITE]: %s -> %s (大小: %s)\n", w.title, filepath.Base(w.GetFileName(title)), humanize.Bytes(uint64(len(ss.Get()))))
	return nil
}

type SafeString struct {
	sync.RWMutex
	s string
}

func (ss *SafeString) Set(s string) {
	ss.Lock()
	defer ss.Unlock()
	ss.s = s
}
func (ss *SafeString) Replace(o, n string) {
	ss.Lock()
	defer ss.Unlock()
	ss.s = strings.ReplaceAll(ss.s, o, n)
}

func (ss *SafeString) Get() string {
	ss.RLock()
	defer ss.RUnlock()
	return ss.s
}

func FindAllImages(md string) (images []string) {
	for _, matches := range imgRegexp.FindAllStringSubmatch(md, -1) {
		if len(matches) == 3 {
			images = append(images, matches[2])
		}
	}
	return
}
