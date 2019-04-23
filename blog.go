package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/russross/blackfriday"
	"github.com/shibukawa/git4go"
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

var (
	lastCommit string
	blogCache  = map[string]string{}
)

func validateCache() {
	commit := getLastCommit()
	if commit != lastCommit {
		blogCache = map[string]string{}
		lastCommit = commit
		tocCache = nil
	}
}

func getBlog(path string) string {
	if !nocache {
		validateCache()
		if blog, ok := blogCache[path]; ok {
			return blog
		}
	}
	blogBytes, err := ioutil.ReadFile(path)
	checkErr(err)
	blog := renderTemplate(BLOG_TPL, parseBlog(blogBytes))
	blogCache[path] = blog
	return blog
}

type blogs []*blogPost

func (b blogs) Len() int           { return len(b) }
func (b blogs) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b blogs) Less(i, j int) bool { return b[i].Date.Before(b[j].Date) }

var tocCache blogs

func getTocData() []*blogPost {
	if !nocache && tocCache != nil {
		return tocCache
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
	sort.Sort(sort.Reverse(blogs(data)))
	return data
}

func getToc() string {
	if !nocache {
		validateCache()
		if toc, ok := blogCache[""]; ok {
			return toc
		}
	}
	data := getTocData()
	toc := renderTemplate(BLOG_TOC, data)
	blogCache[""] = toc
	return toc
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	path := strings.TrimPrefix(r.URL.Path, "/blog")
	if path == "" || path == "/" {
		fmt.Fprint(w, getToc())
		return
	}
	path = BLOG_DIR + path + ".md"
	if !fileExists(path) {
		notfoundHandler(w, r)
		return
	}
	fmt.Fprint(w, getBlog(path))
}

var re = *regexp.MustCompile(`(?i)(title|date|intro|tags|status|toc|position):(.+)`)

type blogPost struct {
	Title          string    `json:"title"`
	Date           time.Time `json:"date"`
	Dateisset      bool      `json:"-"`
	Intro          string    `json:"intro"`
	Tags           []string  `json:"tags"`
	Status         int       `json:"status"`
	Toc            bool      `json:"toc"`
	Position       int       `json:"position"`
	Blog           string    `json:"blog"`
	Url            string    `json:"url"`
	Previous, Next *blogPost
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
	blog, blogText := parseMetadata(blogBytes)
	unsafe := blackfriday.MarkdownCommon(blogText)
	blog.Blog = string(unsafe)
	toc := getTocData()
	for i, b := range toc {
		if blog.Title == b.Title {
			if i > 0 {
				blog.Previous = toc[i-1]
			}
			if i+1 < len(toc) {
				blog.Next = toc[i+1]
			}
			break
		}
	}
	return blog
}
