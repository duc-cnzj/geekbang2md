package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/DuC-cnZj/geekbang2md/cache"

	"github.com/DuC-cnZj/geekbang2md/api"
	"github.com/DuC-cnZj/geekbang2md/read_password"
	"github.com/DuC-cnZj/geekbang2md/zhuanlan"
)

var (
	cookie  string
	noaudio bool
)

func init() {
	flag.StringVar(&cookie, "cookie", "", "-cookie xxxx")
	flag.BoolVar(&noaudio, "noaudio", false, "-noaudio 不下载音频")
}

func main() {
	flag.Parse()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{}, 1)
	go func() {
		select {
		case <-ch:
			done <- struct{}{}
		}
	}()
	go func() {
		var products api.ApiProjectResponse
		var err error
		var phone, password string

		for {
			if cookie == "" {
				if phone == "" || password == "" {
					fmt.Printf("用户名: ")
					fmt.Scanln(&phone)
					password = read_password.ReadPassword("密码: ")
					api.HttpClient.SetPassword(password)
					api.HttpClient.SetPhone(phone)
				}
				if err = api.HttpClient.Login(phone, password, 0); err != nil {
					log.Fatalln(err)
				}
				log.Println("login success")
			} else {
				api.HttpClient.SetHeaders(map[string]string{"Cookie": cookie})
			}
			auth, err := api.HttpClient.UserAuth()
			if err != nil {
				log.Println(err)
				time.Sleep(10 * time.Second)
				continue
			}
			log.Printf("############ %s ############", auth.Data.Nick)
			products, err = api.Products(100)
			if err != nil {
				time.Sleep(10 * time.Second)
				continue
			}
			if products.Code == -1 {
				log.Fatalln("再等等吧, 不让抓了")
			} else {
				break
			}
		}
		for index, product := range products.Data.Products {
			log.Printf("[%d] %s ---%s\n", index+1, product.Title, product.Author.Name)
		}

		var (
			courseID string
			courses  []api.Product
		)
	ASK:
		for {
			courses = nil
			courseID = ""
			fmt.Printf("选择你要爬取的课程(多个用 , 隔开), 直接回车默认全部: \n")
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
		log.Println("############ 爬取的课程 ############")
		for _, cours := range courses {
			log.Printf(cours.Title)
		}
		log.Println("############")

		m := map[int]int{}
		for _, s := range products.Data.List {
			m[s.Pid] = s.Aid
		}
		defer func(t time.Time) { log.Printf("🍌 一共耗时: %s\n", time.Since(t)) }(time.Now())
		wg := sync.WaitGroup{}
		for i := range courses {
			wg.Add(1)
			go func(product *api.Product) {
				defer wg.Done()
				var aid = m[product.ID]
				if aid == 0 && len(product.Column.RecommendArticles) > 0 {
					aid = product.Column.RecommendArticles[0]
				}
				zhuanlan.NewZhuanLan(
					product.Title,
					product.ID,
					aid,
					product.Author.Name,
					product.Article.Count,
					product.Seo.Keywords,
					noaudio,
				).Download()
			}(&courses[i])
		}

		wg.Wait()
		var current, _ = os.Getwd()
		var count int
		filepath.Walk(filepath.Join(current, "geekbang"), func(path string, info fs.FileInfo, err error) error {
			count++
			if info.Mode().IsRegular() && info.Size() < 10 {
				log.Printf("%s 文件为空\n", path)
			}
			return nil
		})
		log.Printf("共计 %d 个文件\n", count)
		log.Println(fmt.Sprintf("缓存位于 %s 目录，可以随意删除", cache.Dir()))
		log.Println("🥭 END")
		done <- struct{}{}
	}()

	<-done
	log.Println("ByeBye")
}
