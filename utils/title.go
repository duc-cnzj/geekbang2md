package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var regexpTitle = regexp.MustCompile(`^(\s*(\d+)\s*|第\d+讲\s)`)
var regexpSpace = regexp.MustCompile(`(\s+)`)

func GetTitle(in string, i int, pad int) string {
	return fmt.Sprintf("%s %s", GetArticleNumber(i, pad), regexpSpace.ReplaceAllString(
		FilterCharacters(regexpTitle.ReplaceAllString(in, "")),
		" "),
	)
}

func GetArticleNumber(i int, pad int) string {
	return fmt.Sprintf("%0*d", pad, i+1)
}

var rd, _ = template.New("").Parse(`
# {{ .Title }}

> 作者: {{ .Author }}
>
> 总数: {{ .Count }}

关键字: {{ .Keywords }}。
`)

func WriteReadmeMD(baseDir, title, author string, count int, keywords []string) error {
	bf := bytes.Buffer{}
	rd.Execute(&bf, map[string]interface{}{
		"Title":    title,
		"Author":   author,
		"Count":    count,
		"Keywords": strings.Join(keywords, ", "),
	})
	file, err := os.OpenFile(filepath.Join(baseDir, "README.md"), os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(bf.Bytes()); err != nil {
		return err
	}
	return nil
}
