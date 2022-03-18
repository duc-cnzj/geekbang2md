package bar

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Interface interface {
	Add()
	AddWithDesc(desc string)
}

type Bar struct {
	sync.Mutex
	title string
	total int
	b     *progressbar.ProgressBar
}

func NewBar(title string, total int) *Bar {
	return &Bar{title: title, total: total, b: progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(60),
		progressbar.OptionSetDescription("[Download]: "+title),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)}
}

func (b *Bar) Add() {
	b.Lock()
	defer b.Unlock()
	b.b.Add(1)
}

func (b *Bar) AddWithDesc(desc string) {
	b.Lock()
	defer b.Unlock()
	b.b.Describe(fmt.Sprintf("%s: (%s)", b.title, desc))
	b.b.Add(1)
	time.Sleep(100 * time.Millisecond)
	if b.b.IsFinished() {
		log.Println()
	}
}
