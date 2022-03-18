package bar

import (
	"github.com/schollz/progressbar/v3"
)

type Interface interface {
	Add()
}

type Bar struct {
	title string
	total int
	b     *progressbar.ProgressBar
}

func NewBar(title string, total int) *Bar {
	return &Bar{title: title, total: total, b: progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(60),
		progressbar.OptionSetDescription("[Download]: "+title),
	)}
}

func (b *Bar) Add() {
	b.b.Add(1)
}
