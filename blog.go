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
	"time"

	"github.com/russross/blackfriday"
	"github.com/shibukawa/git4go"
	"github.com/valyala/fasthttp"
)

const (
	BLOG_DIR = "./blog/"
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
var blogCache = map[string]string{}

func validateCache() {
	commit := getLastCommit()
	if commit != lastCommit {
		blogCache = map[string]string{}
		lastCommit = commit
		tocCache = nil
	}
}

func getBlog(path string) string {
	// validateCache()
	// if blog, ok := blogCache[path]; ok {
	// 	return blog
	// }
	blogBytes, err := ioutil.ReadFile(path)
	checkErr(err)
	blog := renderTemplate(BLOG_TPL, parseBlog(blogBytes))
	blogCache[path] = blog
	return blog
}

var tocCache []*blogPost

func getTocData() []*blogPost {
	// if tocCache != nil {
	// 	return tocCache
	// }
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
	return data
}

func getToc() string {
	// validateCache()
	// if toc, ok := blogCache[""]; ok {
	// 	return toc
	// }
	data := getTocData()
	toc := renderTemplate(BLOG_TOC, data)
	blogCache[""] = toc
	return toc
}

func blogHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html; charset=UTF-8")
	path := strings.TrimPrefix(string(ctx.Path()), "/blog")
	if path == "" || path == "/" {
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
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	Dateisset bool      `json:"-"`
	Intro     string    `json:"intro"`
	Tags      []string  `json:"tags"`
	Status    int       `json:"status"`
	Toc       bool      `json:"toc"`
	Position  int       `json:"position"`
	Blog      string    `json:"blog"`
	Url       string    `json:"url"`
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
			ret.Dateisset = true
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
