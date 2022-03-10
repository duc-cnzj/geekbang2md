package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	sf2 "github.com/DuC-cnZj/geekbang2md/sf"
	"github.com/DuC-cnZj/geekbang2md/waiter"
	"golang.org/x/time/rate"
)

var sf sf2.Group

var HttpClient = newClient()

type client struct {
	mu      sync.RWMutex
	cookies []*http.Cookie
	headers map[string]string

	c               *http.Client
	rt              *waiter.Waiter
	phone, password string
}

func newClient() *client {
	return &client{c: &http.Client{}, rt: waiter.NewWaiter(rate.Every(5*time.Second), 10), headers: map[string]string{}}
}

func (c *client) SetPhone(phone string) {
	c.phone = phone
}
func (c *client) SetPassword(pwd string) {
	c.password = pwd
}

func (c *client) Get(url string, direct bool) (resp *http.Response, err error) {
	r, _ := http.NewRequest("GET", url, nil)
	c.addHeaders(r)

	var do *http.Response
	if direct {
		do, err = c.c.Do(r)
	} else {
		do, err = c.Do(r)
	}
	if err != nil {
		return nil, err
	}
	do, err = c.handleError(do, false)
	if err != nil {
		return nil, err
	}
	var reader io.ReadCloser
	switch do.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(do.Body)
		if err != nil {
			do.Body.Close()
			return nil, err
		}
		do.Body = reader
	default:
	}

	return do, err
}

func (c *client) Post(url string, data interface{}, direct bool) (resp *http.Response, err error) {
	var body io.Reader
	switch d := data.(type) {
	case string:
		body = strings.NewReader(d)
	default:
		marshal, _ := json.Marshal(data)
		body = bytes.NewReader(marshal)
	}
	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	c.addHeaders(r)

	var do *http.Response
	if direct {
		do, err = c.c.Do(r)
	} else {
		do, err = c.Do(r)
	}
	if err != nil {
		return nil, err
	}
	do, err = c.handleError(do, direct)
	if err != nil {
		return nil, err
	}
	var reader io.ReadCloser
	switch do.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(do.Body)
		if err != nil {
			do.Body.Close()
			return nil, err
		}
		do.Body = reader
	default:
	}
	all, _ := io.ReadAll(do.Body)
	e := &GKError{}
	json.Unmarshal(all, &e)
	if e != nil && e.Error.Code < 0 {
		do.Body.Close()
		return nil, errors.New(fmt.Sprintf("%v %d", e.Error.Msg, e.Error.Code))
	}
	do.Body = io.NopCloser(bytes.NewBuffer(all))

	return do, err
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	c.rt.Wait(context.TODO())
	//log.Println("called: ", req.URL)
	var res *http.Response
	var err error
	//log.Printf("call: %s", req.URL)
	res, err = c.c.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (c *client) Login(cellphone, password string) (*AuthInfo, error) {
	c.SetCookies(nil)
	u, err, shared := sf.Do("login", func() (interface{}, error) {
		var user *AuthInfo
		post, err := c.Post("https://account.geekbang.org/account/ticket/login", map[string]interface{}{
			"country":   86,
			"cellphone": cellphone,
			"password":  password,
			"captcha":   "",
			"remember":  1,
			"platform":  3,
			"appid":     1,
			"source":    "",
		}, true)
		if err != nil {
			return nil, err
		}
		defer post.Body.Close()
		if err != nil {
			return nil, err
		}
		c.SetCookies(post.Cookies())
		ti, err := c.Time()
		if err != nil {
			return nil, err
		}
		if user, err = c.UserAuth(ti.Data * 1000); err != nil {
			return nil, err
		}
		log.Println("重新登录成功")
		return user, nil
	})
	if shared {
		log.Println("login request shared.")
	}
	if err != nil {
		return nil, err
	}
	return u.(*AuthInfo), nil
}
func (c *client) Token(token string) error {
	res, err := c.Post("https://account.infoq.cn/account/ticket/token", map[string]interface{}{
		"token": token,
	}, true)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	c.SetCookies(res.Cookies())
	return nil
}

type TimeResponse struct {
	Data int `json:"data"`
	Code int `json:"code"`
}

func (c *client) Time() (*TimeResponse, error) {
	var r *TimeResponse
	res, err := c.Get("https://time.geekbang.org/serv/v1/time", true)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	c.SetCookies(res.Cookies())
	return r, nil
}

type AuthInfo struct {
	Error []interface{} `json:"error"`
	Extra []interface{} `json:"extra"`
	Data  struct {
		Euid        string `json:"euid"`
		Usersubtype int    `json:"usersubtype"`
		Avatar      string `json:"avatar"`
		Usertype    int    `json:"usertype"`
		Cert        int    `json:"cert"`
		Cellphone   string `json:"cellphone"`
		UID         int    `json:"uid"`
		Medalid     int    `json:"medalid"`
		Nick        string `json:"nick"`
		Appid       int    `json:"appid"`
		Ctime       string `json:"ctime"`
		Student     int    `json:"student"`
	} `json:"data"`
	Code int `json:"code"`
}

func (c *client) UserAuth(t int) (*AuthInfo, error) {
	res, err := c.Get("https://account.geekbang.org/serv/v1/user/auth?t="+strconv.Itoa(t), true)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	info := &AuthInfo{}
	json.NewDecoder(res.Body).Decode(&info)
	if info.Code != 0 {
		return nil, errors.New(fmt.Sprintf("%v %d", info.Error, info.Code))
	}
	c.SetCookies(res.Cookies())
	return info, nil
}

func (c *client) SetCookies(cookies []*http.Cookie) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cookies = append(c.cookies, cookies...)
}

func (c *client) SetHeaders(m map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.headers = m
}

type GKError struct {
	Error struct {
		Msg  string `json:"msg"`
		Code int    `json:"code"`
	} `json:"error"`
}

func (c *client) addHeaders(r *http.Request) {
	r.Header.Add("Accept-Encoding", "gzip")
	r.Header.Add("Accept", "application/json, text/plain, */*")
	r.Header.Add("Accept-Encoding", "gzip, deflate, br")
	r.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	r.Header.Add("Connection", "keep-alive")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Host", "time.geekbang.org")
	r.Header.Add("Origin", "https://time.geekbang.org")
	r.Header.Add("Referer", "https://time.geekbang.org/dashboard/course")
	r.Header.Add("Sec-Fetch-Dest", "empty")
	r.Header.Add("Sec-Fetch-Mode", "cors")
	r.Header.Add("Sec-Fetch-Site", "same-origin")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36")
	func() {
		c.mu.RLock()
		defer c.mu.RUnlock()
		for i := range c.cookies {
			r.AddCookie(c.cookies[i])
		}
		for k, v := range c.headers {
			r.Header.Add(k, v)
		}
	}()
}

func (c *client) handleError(do *http.Response, direct bool) (*http.Response, error) {
	if do.StatusCode == 451 || do.StatusCode == 452 {
		defer do.Body.Close()
		if !direct {
			c.rt.Stw()
			time.Sleep(20 * time.Second)
			if _, err := c.Login(c.phone, c.password); err != nil {
				log.Fatalln("login err: ", err)
			}
			c.rt.Restart()
		}
		return nil, errors.New("geekbang 451: 请求太频繁了，再等等吧，程序虽然能继续运行，但还是建议你过会儿再抓")
	}
	if do.StatusCode > 400 {
		defer do.Body.Close()
		all, _ := io.ReadAll(do.Body)
		return nil, errors.New(fmt.Sprintf("%d %v", do.StatusCode, all))
	}
	return do, nil
}
