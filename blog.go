package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/russross/blackfriday"
	"github.com/valyala/fasthttp"
)

const (
	BLOG_DIR = "/home/tso/blog/"
	BLOG_TPL = `<!doctype html>
<html>
	<body>
		<h1>{{.Title}}</h1>
		<h2>{{.Date}}</h2>
		{{.Blog}}
		<hr>
		<address>~ uguu</address>
	</body>
</html>`
)

func blogHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")
	path := strings.TrimPrefix(string(ctx.Path()), "/blog/")
	path = BLOG_DIR + path + ".md"
	if !fileExists(path) {
		notfoundHandler(ctx)
		return
	}
	blogBytes, err := ioutil.ReadFile(path)
	checkErr(err)
	t := template.Must(template.New("blog").Parse(BLOG_TPL))
	err = t.Execute(ctx, parseBlog(blogBytes))
	checkErr(err)
}

var re = *regexp.MustCompile(`(?i)(title|date|intro|tags|status|toc|position):(.+)`)

type blogPost struct {
	Title    string
	Date     time.Time
	Intro    string
	Tags     []string
	Status   int
	Toc      bool
	Position int
	Blog     string
}

func parseBlog(blogBytes []byte) *blogPost {
	ret := &blogPost{}
	buf := bytes.NewBuffer(blogBytes)
	for {
		ln, _ := buf.ReadString('\n')
		if ln == "" || !re.MatchString(ln) {
			break
		}
		m := re.FindStringSubmatch(ln)
		field, value := m[1], m[2]
		value = strings.TrimSpace(value)
		switch strings.ToLower(field) {
		case "title":
			ret.Title = value
		case "date":
			var err error
			ret.Date, err = time.Parse("Mon 02 Jan 2006 15:04:05 PM MST", value)
			checkErr(err)
		case "intro":
			ret.Intro = value
		case "tags":
			ret.Tags = strings.Split(value, " ")
		case "status":
			if value == "public" {
				ret.Status = 0
			} else {
				ret.Status = 1
			}
		case "toc":
			if value == "yes" {
				ret.Toc = true
			} else {
				ret.Toc = false
			}
		case "position":
			ret.Position, _ = strconv.Atoi(value)
		}
	}
	unsafe := blackfriday.MarkdownCommon(buf.Bytes())
	ret.Blog = string(unsafe)
	fmt.Printf("%#v\n", ret)
	return ret
}
