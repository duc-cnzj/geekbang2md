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

func GetTitle(in string, i int, pad int) string {
	title := regexpTitle.ReplaceAllString(in, "")
	return fmt.Sprintf("%0*d %s", pad, i+1, title)
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
	if _, err := file.Write([]byte(bf.String())); err != nil {
		return err
	}
	return nil
}
