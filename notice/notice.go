package notice

import (
	"fmt"
	"log"
	"sync"
)

var n = &notice{courses: make([]fmt.Stringer, 0)}

type warning string

func (w warning) String() string {
	return fmt.Sprintf("❗️ %s", string(w))
}

type course struct {
	title    string
	author   string
	warning  string
	solution string
	ctype    string
}

func (c *course) String() string {
	return fmt.Sprintf("❗️ %s课程 '%s - %s', %s, 解决方案: %s", c.ctype, c.title, c.author, c.warning, c.solution)
}

type notice struct {
	sync.Mutex
	courses []fmt.Stringer
}

func (no *notice) add(c fmt.Stringer) {
	no.Lock()
	defer no.Unlock()
	n.courses = append(n.courses, c)
}

func (no *notice) show() {
	no.Lock()
	defer no.Unlock()
	for _, c := range n.courses {
		log.Println(c)
	}
}

func Warning(w string) {
	n.add(warning(w))
}

func CourseWarning(title, author, warning, solution, ctype string) {
	n.add(&course{
		title:    title,
		author:   author,
		warning:  warning,
		solution: solution,
		ctype:    ctype,
	})
}

func ShowWarnings() {
	n.show()
	log.Println()
}
