package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/DuC-cnZj/geekbang2md/api"
	"github.com/DuC-cnZj/geekbang2md/read_password"
	"github.com/DuC-cnZj/geekbang2md/zhuanlan"
)

var (
	cookie  string
	noaudio bool
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
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
			products, err = api.Products()
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
		m := map[int]int{}
		for _, s := range products.Data.List {
			m[s.Pid] = s.Aid
		}
		wg := sync.WaitGroup{}
		for i := range products.Data.Products {
			wg.Add(1)
			go func(product *api.Product) {
				defer wg.Done()
				log.Printf("%s ---%s\n", product.Title, product.Author.Name)
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
			}(&products.Data.Products[i])
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
		log.Println("🥭 END")
		done <- struct{}{}
	}()

	<-done
	log.Println("ByeBye")
}
