package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/DuC-cnZj/geekbang2md/api"
	"github.com/DuC-cnZj/geekbang2md/cache"
	"github.com/DuC-cnZj/geekbang2md/constant"
	"github.com/DuC-cnZj/geekbang2md/utils"
	"github.com/DuC-cnZj/geekbang2md/video"
	"github.com/DuC-cnZj/geekbang2md/zhuanlan"
)

var (
	dir          string
	cookie       string
	noaudio      bool
	downloadType string
	hack         bool
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&cookie, "cookie", "", "-cookie xxxx")
	flag.BoolVar(&hack, "hack", false, "-hack 获取全部视频，不管你有没有")
	flag.BoolVar(&noaudio, "noaudio", true, "-noaudio 不下载音频")
	flag.StringVar(&dir, "dir", constant.TempDir, fmt.Sprintf("-dir /tmp 下载目录, 默认使用临时目录: '%s'", constant.TempDir))
	flag.StringVar(&downloadType, "type", "", "-type zhuanlan/video 下载类型，不指定则默认全部类型")
}

func main() {
	flag.Parse()
	validateType()

	dir = filepath.Join(dir, "geekbang")
	cache.Init(dir)
	zhuanlan.Init(dir)
	video.Init(dir)

	done := systemSignal()
	go func() {
		var err error
		var phone, password string

		if cookie != "" {
			api.HttpClient.SetHeaders(map[string]string{"Cookie": cookie})
			ti, err := api.HttpClient.Time()
			if err != nil {
				log.Fatalln(err)
			}
			if u, err := api.HttpClient.UserAuth(ti.Data * 1000); err == nil {
				log.Printf("############ %s ############", u.Data.Nick)
			} else {
				log.Fatalln(err)
			}
		} else {
			if phone == "" || password == "" {
				fmt.Printf("用户名: ")
				fmt.Scanln(&phone)
				password = utils.ReadPassword("密码: ")
				api.HttpClient.SetPassword(password)
				api.HttpClient.SetPhone(phone)
			}
			if u, err := api.HttpClient.Login(phone, password); err != nil {
				log.Fatalln(err)
			} else {
				log.Printf("############ %s ############", u.Data.Nick)
			}
		}
		var products api.ProductList
		ptype := api.ProductTypeAll

		switch downloadType {
		case "zhuanlan":
			ptype = api.ProductTypeZhuanlan
		case "video":
			ptype = api.ProductTypeVideo
		}

		if hack {
			products, err = all(ptype)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			products, err = api.AllProducts(ptype)
		}
		if err != nil {
			log.Fatalln("获取课程失败", err)
		}
		courses := prompt(products)
		defer func(t time.Time) { log.Printf("🍌 一共耗时: %s\n", time.Since(t)) }(time.Now())

		for i := range courses {
			var product = &courses[i]
			log.Printf("开始爬取: <%s>\n", product.Title)

			switch product.Type {
			case api.ProductTypeVideo:
				video.NewVideo(
					product.Title,
					product.ID,
					product.Author.Name,
					product.Article.Count,
					product.Seo.Keywords,
				).Download()
			case api.ProductTypeZhuanlan:
				zhuanlan.NewZhuanLan(
					product.Title,
					product.ID,
					product.Author.Name,
					product.Article.Count,
					product.Seo.Keywords,
					noaudio,
				).Download()
			default:
				log.Printf("未知类型, %s\n", product.Type)
			}
		}

		var (
			count     int
			totalSize int64
			cacheSize int64
		)
		filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			count++
			if info.Mode().IsRegular() {
				if strings.HasPrefix(path, cache.Dir()) {
					cacheSize += info.Size()
				}
				if info.Size() < 10 {
					log.Printf("%s 文件为空\n", path)
				}
				totalSize += info.Size()
			}
			return nil
		})
		log.Printf("共计 %d 个文件\n", count)
		log.Printf("🍓 markdown 目录位于: %s, 大小是 %s\n", dir, humanize.Bytes(uint64(totalSize)))
		log.Printf("🥡 缓存目录, 请手动删除: %s, 大小是 %s\n", cache.Dir(), humanize.Bytes(uint64(cacheSize)))
		log.Println("🥭 END")
		done <- struct{}{}
	}()

	<-done
	log.Println("\nByeBye")
}

func all(ptype api.PType) (api.ProductList, error) {
	var products api.ProductList
	skus, err := api.Skus(ptype)
	if err != nil {
		return nil, err
	}
	var chunks [][]string
	var start, end int = 0, 10
	var hasMore bool = true
	for hasMore {
		if len(skus.Data.List) <= end {
			end = len(skus.Data.List)
			hasMore = false
		}
		datas := skus.Data.List[start:end]
		var ids []string
		for _, data := range datas {
			ids = append(ids, strconv.Itoa(data.ColumnSku))
		}
		chunks = append(chunks, ids)
		if hasMore {
			start += 10
			end += 10
		}
	}
	for _, chunk := range chunks {
		infos, err := api.Infos(chunk)
		if err != nil {
			return nil, err
		}
		for _, article := range infos.Data.Infos {
			products = append(products, api.Product{
				ID:       article.ID,
				Type:     article.Type,
				Title:    article.Title,
				Subtitle: article.Subtitle,
				Author: struct {
					Name      string `json:"name"`
					Intro     string `json:"intro"`
					Info      string `json:"info"`
					Avatar    string `json:"avatar"`
					BriefHTML string `json:"brief_html"`
					Brief     string `json:"brief"`
				}{
					Name: article.Author.Name,
				},
				Article: struct {
					ID                int    `json:"id"`
					Count             int    `json:"count"`
					CountReq          int    `json:"count_req"`
					CountPub          int    `json:"count_pub"`
					TotalLength       int    `json:"total_length"`
					FirstArticleID    int    `json:"first_article_id"`
					FirstArticleTitle string `json:"first_article_title"`
				}{
					ID:                article.Article.ID,
					Count:             article.Article.Count,
					CountReq:          article.Article.CountReq,
					CountPub:          article.Article.CountPub,
					TotalLength:       article.Article.TotalLength,
					FirstArticleID:    article.Article.FirstArticleID,
					FirstArticleTitle: article.Article.FirstArticleTitle,
				},
				Seo: struct {
					Keywords []string `json:"keywords"`
				}{
					Keywords: article.Seo.Keywords,
				},
			})
		}
	}
	return products, nil
}

func validateType() {
	if downloadType != "" && downloadType != "zhuanlan" && downloadType != "video" {
		log.Fatalf("type 参数校验失败, '%s' \n", downloadType)
	}
}

func prompt(products api.ProductList) []api.Product {
	sort.Sort(products)
	for index, product := range products {
		var ptypename string
		switch product.Type {
		case api.ProductTypeZhuanlan:
			ptypename = "专栏"
		case api.ProductTypeVideo:
			ptypename = "视频"

		}
		log.Printf("[%d] (%s) %s --- %s\n", index+1, ptypename, product.Title, product.Author.Name)
	}

	var (
		courseID string
		courses  []api.Product
	)
ASK:
	for {
		courses = nil
		courseID = ""
		fmt.Printf("🍎 下载的目录是: '%s', 选择你要爬取的课程(多个用 , 隔开), 直接回车默认全部: \n", dir)
		fmt.Printf("> ")
		fmt.Scanln(&courseID)
		if courseID == "" {
			courses = products
			break
		}
		split := strings.Split(courseID, ",")
		for _, s := range split {
			id, err := strconv.Atoi(s)
			if err != nil || id > len(products) || id < 1 {
				log.Printf("非法课程 id %v !\n", s)
				continue ASK
			}
			courses = append(courses, products[id-1])
		}
		break
	}
	return courses
}

func systemSignal() chan struct{} {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{}, 1)
	go func() {
		select {
		case <-ch:
			done <- struct{}{}
		}
	}()
	return done
}
