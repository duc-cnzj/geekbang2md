package zhuanlan

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/dustin/go-humanize"

	"github.com/DuC-cnZj/geekbang2md/api"
	"github.com/DuC-cnZj/geekbang2md/image"
	"github.com/DuC-cnZj/geekbang2md/markdown"
)

type ZhuanLan struct {
	noaudio  bool
	title    string
	id       int
	author   string
	count    int
	keywords []string
	aid      int

	imageManager *image.Manager
	mdWriter     *markdown.MDWriter
}

var current, _ = os.Getwd()

var imageManager = image.NewManager(filepath.Join(current, "geekbang", "images"))

func NewZhuanLan(title string, id, aid int, author string, count int, keywords []string, noaudio bool) *ZhuanLan {
	mdWriter := markdown.NewMDWriter(filepath.Join(current, "geekbang", title), title, imageManager)
	return &ZhuanLan{noaudio: noaudio, title: title, id: id, aid: aid, author: author, count: count, keywords: keywords, imageManager: imageManager, mdWriter: mdWriter}
}

var rd, _ = template.New("").Parse(`
# {{ .Title }}

> author: {{ .Author }}
>
> count: {{ .Count }}

keywords: {{ .Keywords }}。
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
	article, err := api.Article(strconv.Itoa(zl.aid))
	if err != nil {
		log.Println(err, zl.aid)
		return err
	}
	articles, err := api.Articles(article.Data.Cid)
	if err != nil {
		log.Println(err)
	}
	wg := sync.WaitGroup{}
	for i := range articles.Data.List {
		wg.Add(1)
		go func(s *api.ArticlesResponseItem) {
			defer wg.Done()
			if zl.mdWriter.FileExists(s.ArticleTitle) {
				//log.Println("[SKIP]: ", s.ArticleTitle)
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
				if err := zl.mdWriter.WriteFile(s.AudioDownloadURL, s.AudioDubber, humanize.Bytes(uint64(s.AudioSize)), s.AudioTime, s.ArticleTitle, response.Data.ArticleContent); err != nil {
					log.Println(err)
				}
			}
		}(articles.Data.List[i])
	}

	wg.Wait()
	return nil
}
