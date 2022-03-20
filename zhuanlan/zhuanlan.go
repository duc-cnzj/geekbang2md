package zhuanlan

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/duc-cnzj/geekbang2md/api"
	"github.com/duc-cnzj/geekbang2md/bar"
	"github.com/duc-cnzj/geekbang2md/image"
	"github.com/duc-cnzj/geekbang2md/utils"
)

type ZhuanLan struct {
	audio    bool
	title    string
	id       int
	author   string
	count    int
	keywords []string

	imageManager *image.Manager
	mdWriter     *MDWriter
}

var baseDir string

func Init(d string) {
	baseDir = d
}

func NewZhuanLan(title string, id int, author string, count int, keywords []string, audio bool) *ZhuanLan {
	dir := filepath.Join(baseDir, utils.FilterCharacters(title))
	imageManager := image.NewManager(filepath.Join(dir, "images"))

	mdWriter := NewMDWriter(dir, title, imageManager)
	return &ZhuanLan{audio: audio, title: title, id: id, author: author, count: count, keywords: keywords, imageManager: imageManager, mdWriter: mdWriter}
}

func (zl *ZhuanLan) Download() error {
	utils.WriteReadmeMD(zl.mdWriter.baseDir, zl.title, zl.author, zl.count, zl.keywords)
	articles, err := api.Articles(zl.id)
	if err != nil {
		return err
	}
	var pad int = 2
	if zl.count > 100 {
		pad = 3
	}
	currentCount := len(articles.Data.List)
	b := bar.NewBar(zl.title, currentCount)
	wg := sync.WaitGroup{}
	r := NewZlResults()
	for i := range articles.Data.List {
		wg.Add(1)
		go func(s *api.ArticlesResponseItem, i int) {
			defer b.Add()
			defer wg.Done()
			t := utils.GetTitle(s.ArticleTitle, i, pad)
			articleNumber := utils.GetArticleNumber(i, pad)
			if !zl.audio {
				s.AudioDownloadURL = ""
			}
			if info, path, exists := zl.mdWriter.FileExists(t); exists {
				skip := true
				file, _ := os.ReadFile(path)
				images := FindAllImages(string(file))
				if s.AudioDownloadURL != "" {
					images = append(images, s.AudioDownloadURL)
				}
				if len(images) > 0 {
					for _, imageUrl := range images {
						localPath, err := zl.imageManager.FullLocalPath(imageUrl, articleNumber)
						if err != nil {
							skip = false
							break
						}
						stat, err := os.Stat(localPath)
						if err != nil {
							skip = false
							break
						}
						if stat.Size() < 10 {
							skip = false
							break
						}
					}
				}
				if skip {
					r.Add(i, fmt.Sprintf("[SKIP]: %s (大小: %s)", filepath.Base(info.Name()), utils.Bytes(uint64(info.Size()))))
					return
				}
			}
			response, err := api.Article(strconv.Itoa(s.ID))
			if err != nil {
				log.Println(err, response.Code)
				return
			}

			if len(response.Data.ArticleContent) > 0 {
				if reason, err := zl.mdWriter.WriteFile(articleNumber, s.AudioDownloadURL, s.AudioDubber, utils.Bytes(uint64(s.AudioSize)), s.AudioTime, t, response.Data.ArticleContent); err != nil {
					r.Add(i, fmt.Sprintf("[下载出错] %s: '%v'", t, err.Error()))
				} else {
					r.Add(i, reason)
				}
			}
		}(articles.Data.List[i], i)
	}

	wg.Wait()
	if zl.count > currentCount {
		api.DeleteArticlesCache(zl.id)
	}
	time.Sleep(300 * time.Millisecond)
	r.Print()
	return nil
}

type ZlResult struct {
	id   int
	info string
}
type SortZlResults []*ZlResult

func (s SortZlResults) Len() int {
	return len(s)
}

func (s SortZlResults) Less(i, j int) bool {
	return s[i].id < s[j].id
}

func (s SortZlResults) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ZlResults struct {
	sync.Mutex
	results []*ZlResult
}

func NewZlResults() *ZlResults {
	return &ZlResults{results: make([]*ZlResult, 0)}
}

func (zlr *ZlResults) Add(id int, result string) {
	zlr.Lock()
	defer zlr.Unlock()
	zlr.results = append(zlr.results, &ZlResult{
		id:   id,
		info: result,
	})
}

func (zlr *ZlResults) Print() {
	zlr.Lock()
	defer zlr.Unlock()
	sort.Sort(SortZlResults(zlr.results))
	log.Println()
	for _, result := range zlr.results {
		log.Println(result.info)
	}
}

var regexpTitle = regexp.MustCompile(`^(\s*(\d+)\s*|第\d+讲\s)`)

func getTitle(s *api.ArticlesResponseItem, i int, pad int) string {
	title := regexpTitle.ReplaceAllString(s.ArticleTitle, "")
	return fmt.Sprintf("%0*d %s", pad, i+1, title)
}
