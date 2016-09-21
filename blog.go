package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/russross/blackfriday"
	"github.com/shibukawa/git4go"
	"github.com/valyala/fasthttp"
)

const (
	BLOG_DIR = "/home/tso/blog/"
	BLOG_TPL = "./templates/blog_entry.html"
	BLOG_TOC = "./templates/blog_index.html"
)

// grab git HEAD commit hash + precompile table of contents
// cache compiled blog posts map[string]string{"/path/to/blog.md": "<html>..."}
// invalidate cache on HEAD commit change, 20 GOTO 10

func getLastCommit() string {
	repo, err := git4go.OpenRepository(BLOG_DIR)
	checkErr(err)
	head, err := repo.Head()
	checkErr(err)
	return head.Target().String()
}

var lastCommit string
var blogCache map[string]string

func validateCache() {
	commit := getLastCommit()
	if commit != lastCommit {
		blogCache = map[string]string{}
		lastCommit = commit
	}
}

func getBlog(path string) string {
	validateCache()
	if blog, ok := blogCache[path]; ok {
		return blog
	}
	blogBytes, err := ioutil.ReadFile(path)
	checkErr(err)
	t, err := template.ParseFiles(BLOG_TPL)
	checkErr(err)
	buf := new(bytes.Buffer)
	checkErr(t.Execute(buf, parseBlog(blogBytes)))
	blog := buf.String()
	blogCache[path] = blog
	return blog
}

func getToc() string {
	validateCache()
	if toc, ok := blogCache[""]; ok {
		return toc
	}
	data := []*blogPost{}
	filepath.Walk(BLOG_DIR, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			blogBytes, err := ioutil.ReadFile(path)
			checkErr(err)
			blogPost, _ := parseMetadata(blogBytes)
			blogPost.Url = strings.TrimSuffix(strings.TrimPrefix(info.Name(), BLOG_DIR), ".md")
			data = append(data, blogPost)
		}
		return nil
	})
	t, err := template.ParseFiles(BLOG_TOC)
	checkErr(err)
	buf := new(bytes.Buffer)
	checkErr(t.Execute(buf, data))
	toc := buf.String()
	blogCache[""] = toc
	return toc
}

func blogHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html; charset=UTF-8")
	path := strings.TrimPrefix(string(ctx.Path()), "/blog/")
	if path == "" {
		fmt.Fprint(ctx, getToc())
		return
	}
	path = BLOG_DIR + path + ".md"
	if !fileExists(path) {
		notfoundHandler(ctx)
		return
	}
	fmt.Fprint(ctx, getBlog(path))
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
	Url      string
}

func parseMetadata(blogBytes []byte) (*blogPost, []byte) {
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
	return ret, buf.Bytes()
}
func parseBlog(blogBytes []byte) *blogPost {
	ret, blog := parseMetadata(blogBytes)
	unsafe := blackfriday.MarkdownCommon(blog)
	ret.Blog = string(unsafe)
	return ret
}
