package bar

import (
	"fmt"

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
	total += 1
	runes := []rune(title)
	if len(runes) > 20 {
		runes = append(runes[0:17], []rune("...")...)
	}
	b := &Bar{title: title, total: total, b: progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(20),
		progressbar.OptionSetDescription(fmt.Sprintf("[Download]: %-20s", string(runes))),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "#",
			SaucerHead:    "#",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)}
	b.b.Add(1)
	return b
}

func (b *Bar) Add() {
	b.b.Add(1)
}
