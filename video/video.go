package video

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"

	"github.com/DuC-cnZj/geekbang2md/utils"

	"github.com/dustin/go-humanize"

	"github.com/DuC-cnZj/geekbang2md/api"
	"github.com/DuC-cnZj/geekbang2md/constant"
	"github.com/DuC-cnZj/geekbang2md/waiter"

	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Video struct {
	sync.RWMutex
	baseDir string
	waiter  waiter.Interface

	cid   int
	title string

	author   string
	count    int
	keywords []string
}

var baseDir string

func Init(d string) {
	baseDir = filepath.Join(d, "videos")
}

var uregex = regexp.MustCompile(`URI="(.*?)"`)

func NewVideo(title string, id int, author string, count int, keywords []string) *Video {
	d := filepath.Join(baseDir, title)
	os.MkdirAll(d, 0755)
	return &Video{
		title:    title,
		baseDir:  d,
		waiter:   waiter.NewSigWaiter(constant.VideoDownloadParallelNum),
		cid:      id,
		author:   author,
		count:    count,
		keywords: keywords,
	}
}

var tsFileRegexp = regexp.MustCompile(`\w+.*?-(\d+)\.ts`)

type Seg struct {
	id      int
	path    string
	fullUrl string
}

type Segs []*Seg

func (s Segs) Len() int {
	return len(s)
}

func (s Segs) Less(i, j int) bool {
	return s[i].id < s[j].id
}

func (s Segs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (v *Video) DownloadPath(name string) string {
	return filepath.Join(v.baseDir, utils.FilterCharacters(name))
}

func (v *Video) SegDownloadPath(name string) string {
	return filepath.Join(v.baseDir, "segs", utils.FilterCharacters(name))
}

func (v *Video) DeleteSegs(segs []*Seg) error {
	for _, seg := range segs {
		os.Remove(seg.path)
	}
	return nil
}

func (v *Video) Download() error {
	articles, err := api.Articles(v.cid)
	if err != nil {
		log.Println(err)
		return err
	}

	wg := sync.WaitGroup{}
	for i := range articles.Data.List {
		wg.Add(1)
		go func(s *api.ArticlesResponseItem) {
			defer wg.Done()
			v.waiter.Wait(context.TODO())
			defer v.waiter.Release()
			article, err := api.Article(strconv.Itoa(s.ID))
			if err != nil {
				return
			}
			marshal, _ := json.Marshal(article.Data.HlsVideos)
			var vi api.Video
			json.Unmarshal(marshal, &vi)
			log.Printf("开始下载: %s", s.ArticleTitle)
			err = download(v.DownloadPath(s.ArticleTitle+".ts"), vi.Hd.URL, v, s)
			if err != nil {
				log.Printf("下载出错: %v\n", err)
			}
		}(articles.Data.List[i])
	}
	wg.Wait()
	p := v.DownloadPath("segs")
	var count int
	filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if d.Type().IsRegular() && strings.HasSuffix(path, ".ts") {
			count++
		}
		return nil
	})
	if count == 0 {
		os.RemoveAll(p)
	}
	return nil
}

func download(path string, u string, v *Video, s *api.ArticlesResponseItem) error {
	var err error
	stat, err := os.Stat(path)
	if err == nil && stat.Size() > 0 {
		return nil
	}
	get, err := api.NewBackoffClient(3).Get(u)
	if err != nil {
		log.Fatalln(err)
	}
	defer get.Body.Close()
	keyAll, err := io.ReadAll(get.Body)
	if err != nil {
		log.Fatalln(err)
	}
	submatch := uregex.FindStringSubmatch(string(keyAll))
	key, _ := api.VideoKey(submatch[1], strconv.Itoa(s.ID))

	parse, _ := url.Parse(u)
	baseUrl := fmt.Sprintf("https://%s/%s/", parse.Host, strings.Split(parse.Path, "/")[1])
	res, err := api.NewBackoffClient(3).Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	all, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	os.MkdirAll(v.DownloadPath("segs"), 0755)
	var items Segs
	stringSubmatch := tsFileRegexp.FindAllStringSubmatch(string(all), -1)
	for _, s := range stringSubmatch {
		id, _ := strconv.Atoi(s[1])
		items = append(items, &Seg{
			id:      id,
			path:    v.SegDownloadPath(s[0]),
			fullUrl: baseUrl + s[0],
		})
	}

	wg := sync.WaitGroup{}
	sigWaiter := waiter.NewSigWaiter(5)
	for i := range items {
		wg.Add(1)
		go func(s *Seg) {
			defer wg.Done()
			sigWaiter.Wait(context.TODO())
			defer sigWaiter.Release()
			st, err := os.Stat(s.path)
			if err == nil && st.Size() > 0 {
				//log.Printf("%s exists", s.path)
				return
			}
			get, err := api.NewBackoffClient(3).Get(s.fullUrl)
			if err != nil {
				return
			}
			defer get.Body.Close()
			readAll, _ := io.ReadAll(get.Body)
			if len(readAll) > 0 {
				//log.Printf("[WRITE]: %s\n", s.path)
				if err := os.WriteFile(s.path, readAll, 0644); err != nil {
					log.Fatalln(err)
				}
			}
		}(items[i])
	}

	wg.Wait()
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	sort.Sort(items)

	for _, item := range items {
		file, err := os.ReadFile(item.path)
		if err != nil {
			return err
		}
		aes128, err := decryptAES128(file, key, make([]byte, 16))
		if err != nil {
			v.DeleteSegs(items)
			log.Printf("[ERROR]: 解码失败: %v\n", err)
			return err
		}
		for j := 0; j < len(aes128); j++ {
			if aes128[j] == syncByte {
				aes128 = aes128[j:]
				break
			}
		}
		if _, err := f.Write(aes128); err != nil {
			return err
		}
	}
	v.DeleteSegs(items)
	info, _ := f.Stat()
	log.Printf("[SUCCESS]: 下载成功 '%s', 大小: '%s'", s.ArticleTitle, humanize.Bytes(uint64(info.Size())))
	return nil
}

const (
	syncByte = uint8(71) //0x47
)

func decryptAES128(crypted, key, iv []byte) (origData []byte, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errors.New(e.(string))
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData = make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs5UnPadding(origData)
	return
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
