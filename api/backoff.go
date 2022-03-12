package api

import (
	"log"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type BackoffClient struct {
	RetryTimes uint64
	c          *http.Client
}

func NewBackoffClient(retryTimes uint64) *BackoffClient {
	return &BackoffClient{RetryTimes: retryTimes, c: &http.Client{}}
}

func (b *BackoffClient) Get(u string) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)
	if err = backoff.Retry(func() (e error) {
		resp, e = b.c.Get(u)
		if e != nil {
			log.Printf("http '%s' err: '%v'  , retry...\n", u, e)
			return e
		}
		return nil
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(3*time.Second), b.RetryTimes)); err != nil {
		return nil, err
	}
	return resp, err
}
