package zhuanlan

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/dustin/go-humanize"

	"github.com/DuC-cnZj/geekbang2md/api"
	"github.com/DuC-cnZj/geekbang2md/image"
	"github.com/DuC-cnZj/geekbang2md/utils"
)

type ZhuanLan struct {
	noaudio  bool
	title    string
	id       int
	author   string
	count    int
	keywords []string

	imageManager *image.Manager
	mdWriter     *MDWriter
}

var baseDir string
var imageManager *image.Manager

func Init(d string) {
	baseDir = d
	imageManager = image.NewManager(filepath.Join(baseDir, "images"))
}

func NewZhuanLan(title string, id int, author string, count int, keywords []string, noaudio bool) *ZhuanLan {
	mdWriter := NewMDWriter(filepath.Join(baseDir, utils.FilterCharacters(title)), title, imageManager)
	return &ZhuanLan{noaudio: noaudio, title: title, id: id, author: author, count: count, keywords: keywords, imageManager: imageManager, mdWriter: mdWriter}
}

var rd, _ = template.New("").Parse(`
# {{ .Title }}

> 作者: {{ .Author }}
>
> 总数: {{ .Count }}

关键字: {{ .Keywords }}。
`)

func (zl *ZhuanLan) Download() error {
	bf := bytes.Buffer{}
	rd.Execute(&bf, map[string]interface{}{
		"Title":    zl.title,
		"Author":   zl.author,
		"Count":    zl.count,
		"Keywords": strings.Join(zl.keywords, ", "),
	})
	zl.mdWriter.WriteReadmeMD(bf.String())
	articles, err := api.Articles(zl.id)
	if err != nil {
		log.Println(err)
		return err
	}
	var pad int = 2
	if zl.count > 100 {
		pad = 3
	}

	wg := sync.WaitGroup{}
	for i := range articles.Data.List {
		wg.Add(1)
		go func(s *api.ArticlesResponseItem, i int) {
			defer wg.Done()
			t := getTitle(s, i, pad)
			if info, exists := zl.mdWriter.FileExists(t); exists {
				log.Printf("[SKIP]: %s -> %s (大小: %s)\n", zl.title, filepath.Base(info.Name()), humanize.Bytes(uint64(info.Size())))
				return
			}
			response, err := api.Article(strconv.Itoa(s.ID))
			if err != nil {
				log.Println(err, response.Code)
				return
			}

			if len(response.Data.ArticleContent) > 0 {
				if zl.noaudio {
					s.AudioDownloadURL = ""
				}
				if err = zl.mdWriter.WriteFile(s.AudioDownloadURL, s.AudioDubber, humanize.Bytes(uint64(s.AudioSize)), s.AudioTime, t, response.Data.ArticleContent); err != nil {
					log.Println(err)
				}
			}
		}(articles.Data.List[i], i)
	}

	wg.Wait()
	return nil
}

var regexpTitle = regexp.MustCompile(`^(\s*(\d+)\s*|第\d+讲\s)`)

func getTitle(s *api.ArticlesResponseItem, i int, pad int) string {
	title := regexpTitle.ReplaceAllString(s.ArticleTitle, "")
	return fmt.Sprintf("%0*d %s", pad, i+1, title)
}
