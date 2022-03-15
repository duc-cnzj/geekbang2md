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
	"sync"
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
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&cookie, "cookie", "", "-cookie xxxx")
	flag.BoolVar(&noaudio, "noaudio", false, "-noaudio 不下载音频")
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

		var products api.ProjectResponse
		ptype := api.ProductTypeAll

		switch downloadType {
		case "zhuanlan":
			ptype = api.ProductTypeZhuanlan
		case "video":
			ptype = api.ProductTypeVideo
		}

		products, err = api.Products(100, ptype)
		if err != nil {
			log.Fatalln("获取课程失败", err)
		}
		if products.Code == -1 {
			log.Fatalln("再等等吧, 不让抓了")
		}
		courses := prompt(products)

		defer func(t time.Time) { log.Printf("🍌 一共耗时: %s\n", time.Since(t)) }(time.Now())

		wg := sync.WaitGroup{}
		for i := range courses {
			wg.Add(1)
			go func(product *api.Product) {
				defer wg.Done()
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
			}(&courses[i])
		}

		wg.Wait()
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
	log.Println("ByeBye")
}

func validateType() {
	if downloadType != "" && downloadType != "zhuanlan" && downloadType != "video" {
		log.Fatalf("type 参数校验失败, '%s' \n", downloadType)
	}
}

func prompt(products api.ProjectResponse) []api.Product {
	sort.Sort(products.Data.Products)
	for index, product := range products.Data.Products {
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
			courses = products.Data.Products
			break
		}
		split := strings.Split(courseID, ",")
		for _, s := range split {
			id, err := strconv.Atoi(s)
			if err != nil || id > len(products.Data.Products) || id < 1 {
				log.Printf("非法课程 id %v !\n", s)
				continue ASK
			}
			courses = append(courses, products.Data.Products[id-1])
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
