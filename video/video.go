package video

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/DuC-cnZj/geekbang2md/api"
	"github.com/DuC-cnZj/geekbang2md/bar"
	"github.com/DuC-cnZj/geekbang2md/constant"
	"github.com/DuC-cnZj/geekbang2md/utils"
	"github.com/DuC-cnZj/geekbang2md/waiter"
)

type Video struct {
	sync.RWMutex
	baseDir string

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
	d := filepath.Join(baseDir, utils.FilterCharacters(title))
	os.MkdirAll(d, 0755)
	return &Video{
		title:    title,
		baseDir:  d,
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

func (v *Video) DeleteSegs(segs ...*Seg) error {
	for _, seg := range segs {
		if err := os.Remove(seg.path); err != nil {
			log.Printf("remove '%s', err: %v", seg.path, err)
		}
	}
	return nil
}

func (v *Video) Download() error {
	utils.WriteReadmeMD(v.baseDir, v.title, v.author, v.count, v.keywords)
	articles, err := api.Articles(v.cid)
	if err != nil {
		return err
	}
	for i := range articles.Data.List {
		func(num int) {
			s := articles.Data.List[i]
			article, err := api.Article(strconv.Itoa(s.ID))
			if err != nil {
				log.Printf("[Download]: article: %s err: %v \n", s.ArticleTitle, err)
				return
			}
			marshal, _ := json.Marshal(article.Data.HlsVideos)
			var vi api.Video
			json.Unmarshal(marshal, &vi)
			if vi.Hd.URL == "" {
				api.DeleteArticleCache(strconv.Itoa(s.ID))
				log.Printf("[ERROR]: 视频: '%s', 下载地址为空！ \n", s.ArticleTitle)
				return
			}
			var pad int = 2
			if v.count > 100 {
				pad = 3
			}
			title := utils.GetTitle(s.ArticleTitle, num, pad)
			for i := 0; i < 3; i++ {
				err = download(v.DownloadPath(title+".ts"), vi.Hd.URL, v, title, strconv.Itoa(s.ID))
				if !errors.Is(err, ErrorRetry) {
					break
				}
				log.Printf("\n[Warning]: 下载出错, 重新下载: '%s', %v\n", title, err)
				time.Sleep(500 * time.Millisecond)
			}
			if err != nil {
				log.Printf("\n下载出错: %v\n", err)
			}
		}(i)
	}
	p := v.DownloadPath("segs")
	if _, err := os.Stat(p); err == nil {
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
	}
	return nil
}

var ErrorRetry = errors.New("retry")

func download(downloadPath string, hdUrl string, v *Video, title string, id string) error {
	var err error
	stat, err := os.Stat(downloadPath)
	if err == nil && stat.Size() > 0 {
		return nil
	}

	parse, err := url.Parse(hdUrl)
	if err != nil {
		return err
	}

	get, err := api.NewBackoffClient(3).Get(hdUrl)
	if err != nil {
		return err
	}
	defer get.Body.Close()
	m3u8, err := io.ReadAll(get.Body)
	if err != nil {
		return err
	}

	baseUrl := fmt.Sprintf("https://%s/%s/", parse.Host, strings.Split(parse.Path, "/")[1])
	os.MkdirAll(v.DownloadPath("segs"), 0755)
	var items Segs
	stringSubmatch := tsFileRegexp.FindAllStringSubmatch(string(m3u8), -1)
	for _, s := range stringSubmatch {
		id, _ := strconv.Atoi(s[1])
		items = append(items, &Seg{
			id:      id,
			path:    v.SegDownloadPath(s[0]),
			fullUrl: baseUrl + s[0],
		})
	}

	wg := sync.WaitGroup{}
	sigWaiter := waiter.NewSigWaiter(constant.VideoDownloadParallelNum)
	var b bar.Interface = bar.NewBar(title, len(items))
	for i := range items {
		wg.Add(1)
		sigWaiter.Wait(context.TODO())
		go func(s *Seg) {
			defer wg.Done()
			defer b.Add()
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
			if get.ContentLength > 0 {
				file, err := os.OpenFile(s.path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					return
				}
				defer file.Close()
				io.Copy(file, bufio.NewReaderSize(get.Body, 1024*1024*10))
			}
		}(items[i])
	}

	wg.Wait()
	sort.Sort(items)
	submatch := uregex.FindStringSubmatch(string(m3u8))
	key, err := api.VideoKey(submatch[1], id)
	if err != nil {
		return err
	}
	if len(key) == 0 {
		api.DeleteArticleCache(id)
		return fmt.Errorf("%w, 当前获取不到解码的 key 值，建议全部下载完成之后再重试。", ErrorRetry)
	}

	f, err := os.OpenFile(downloadPath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, item := range items {
		file, err := os.ReadFile(item.path)
		if err != nil {
			return err
		}
		aes128, err := decryptAES128(file, key, make([]byte, 16))
		if err != nil {
			v.DeleteSegs(items...)
			f.Close()
			os.Remove(downloadPath)
			return fmt.Errorf("[%w]: reason: '%s' path: '%s'", ErrorRetry, err.Error(), item.path)
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
	v.DeleteSegs(items...)
	info, _ := f.Stat()
	log.Printf("\n[SUCCESS]: 下载成功 '%s', 大小: '%s'", title, humanize.Bytes(uint64(info.Size())))
	return nil
}

const (
	syncByte = uint8(71) //0x47
)

func decryptAES128(crypted, key, iv []byte) (origData []byte, err error) {
	defer func() {
		e := recover()
		switch edata := e.(type) {
		case string:
			err = errors.New(fmt.Sprintf("%s, len key: %d", edata, len(key)))
		case error:
			err = fmt.Errorf("%w: key len: %d", edata, len(key))
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: key len: %d", err, len(key))
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
