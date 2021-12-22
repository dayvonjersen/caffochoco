package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type syncMap struct {
	m  map[string]int
	mu sync.Mutex
}

func (s *syncMap) Set(k string, v int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[k] = v
}

func (s *syncMap) Get(k string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.m[k]; ok {
		delete(s.m, k)
		return v, true
	}
	return 0, false
}

func flush(w http.ResponseWriter) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func send(w http.ResponseWriter, msg string) {
	io.WriteString(w, msg)
	flush(w)
	<-time.After(277 * time.Millisecond) // approx speed of 28.8kbps modem
}

func uuid() string {
	lmao := ""
	for i := 0; i < 64; i++ {
		lmao = lmao + fmt.Sprintf("%02x", rand.Intn(255))
	}
	return lmao
}

func head(w http.ResponseWriter, code int, headers map[string]string) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	// HACK(tso):
	// HACK(tso):
	// HACK(tso):
	if code == 0 {
		log.Println("<- THEY GOT A GARBAGE FILE OHNOES")
		return
	}
	log.Println("<-", code)
	w.WriteHeader(code)
}

func file(w http.ResponseWriter, filename string) {
	f, err := os.Open(filename)
	defer f.Close()
	checkErr(err)
	s, err := f.Stat()
	checkErr(err)
	head(w, 300, map[string]string{"Content-Length": fmt.Sprintf("%d", s.Size())})
	io.Copy(w, f)
}

func main() {
	var (
		addr, certFile, keyFile string
		port                    int
	)
	flag.StringVar(
		&addr,
		"addr",
		"",
		"leave blank for 0.0.0.0",
	)
	flag.IntVar(
		&port,
		"port",
		8080,
		"",
	)
	flag.StringVar(
		&certFile,
		"cert",
		"",
		"path to cert",
	)
	flag.StringVar(
		&keyFile,
		"key",
		"",
		"path to key",
	)
	flag.Parse()

	rand.Seed(time.Now().Unix())

	tokens := &syncMap{m: map[string]int{}}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		req(r)
		if r.URL.RawQuery != "" {
			t := strings.Split(r.URL.RawQuery, "=")
			if v, ok := tokens.Get(t[0]); ok {
				switch v {
				case 0:
					head(w, 0, map[string]string{
						"Content-Type":        "application/octet-stream",
						"Content-Disposition": "attachment; filename=\"GARBAGE.zip.001\"",
					})
					http.ServeFile(w, r, "z/GARBAGE.zip.001")
					return
				case 1:
					head(w, 0, map[string]string{
						"Content-Type":        "application/octet-stream",
						"Content-Disposition": "attachment; filename=\"GARBAGE.zip.002\"",
					})
					http.ServeFile(w, r, "z/GARBAGE.zip.002")
					return
				case 2:
					head(w, 0, map[string]string{
						"Content-Type":        "application/octet-stream",
						"Content-Disposition": "attachment; filename=\"GARBAGE.zip.003\"",
					})
					http.ServeFile(w, r, "z/GARBAGE.zip.003")
					return
				case 3:
					head(w, 0, map[string]string{
						"Content-Type":        "application/octet-stream",
						"Content-Disposition": "attachment; filename=\"GARBAGE.zip.004\"",
					})
					http.ServeFile(w, r, "z/GARBAGE.zip.004")
					return
				case 4:
					head(w, 0, map[string]string{
						"Content-Type":        "application/octet-stream",
						"Content-Disposition": "attachment; filename=\"GARBAGE.zip.0xx.tar\"",
					})
					http.ServeFile(w, r, "z/GARBAGE.zip.0xx.tar")
					return
				case 5:
					head(w, 0, map[string]string{
						"Content-Type":        "application/octet-stream",
						"Content-Disposition": "attachment; filename=\"GARBAGE.zip.037\"",
					})
					http.ServeFile(w, r, "z/GARBAGE.zip.037")
					return
				}
				return
			}
			head(w, 417, nil)
			send(w, "What? No, no, no. I was disappointed that you tried.")
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		dir := strings.Split(path, "/")[0]

		if path == "" || strings.HasPrefix(path, "index.") {
			file(w, "dream.html")
			return
		}

		if path == "favicon.ico" {
			file(w, "favicon.ico")
			return
		}
		if path == "SECTOR_MAP.nfo" {
			file(w, "SECTOR_MAP.nfo")
			return
		}

		if fileExists(dir) && (dir == "a" || dir == "v") && fileExists(path) && !isDir(path) {
			http.ServeFile(w, r, path)
			return
		}

		headers := map[string]string{
			"Content-Type":  "text/html",
			"Cache-Control": "no-cache, no-store, must-revalidate",
			"Pragma":        "no-cache",
			"Expires":       "0",
		}
		if r.Method == "POST" {
			r.ParseForm()
			switch path {
			case "0x01":
				phrase := r.FormValue("passphrase")
				penntest := r.FormValue("penntest")
				if phrase == "" && penntest == "" {
					head(w, 401, headers)
					send(w, "You're not even trying. I should &lt;insert politically correct word meaning 'fuck off'&gt;list you")
					return
				}
				if strings.ToLower(phrase) == "old age should burn and rave at close of day" {
					head(w, 202, headers)
					send(w, "WE PRESENT YOU A NEW QUEST.")
					send(w, "<iframe width=560 height=315 src='https://youtube.com/embed/RaW-opIh7GE' title='youch00b data mining EPIC CONTENT tracker :o' frameborder=0></iframe>")
					send(w, "<br>So what arbitrary set of numbers am I looking for? I think I made it pretty obvious that I have literally lost my mind (and I miss it ever so much ;)<br><form action='/0x01' method='post'><input type='number' name='penntest'></form>")
					return
				} else if penntest == "17623654" {
					head(w, 203, headers)
					send(w, "A WinRAR is you.")
					t := uuid()
					tokens.Set(t, 1)
					send(w, `<meta http-equiv='refresh' content="2;URL='/?`+t+`=1'">`)
					return
				} else {
					head(w, 406, headers)
					send(w, "Take a closer look.")
					return
				}
			case "0x03":
				code := r.FormValue("code")
				if code == "" {
					head(w, 400, headers)
					send(w, "no u")
					return
				}
				if code == "&#8805;097483" {
					head(w, 206, headers)
					send(w, "you know you've probably seen this movie as many times as I have")
					t := uuid()
					tokens.Set(t, 3)
					send(w, `<meta http-equiv='refresh' content="2;URL='/?`+t+`=3'">`)
					return
				} else {
					head(w, 428, headers)
					send(w, "just give it up.")
				}
			default:
				head(w, 418, headers)
				send(w, "<pre>here is my handle here is my spout :3")
			}
			return
		}
		switch path {
		// SECTOR 0
		// SECTOR 0
		// SECTOR 0
		case "'sector 0x00'.torrent":
			head(w, 200, headers)
			send(w, "<html><body style='user-select:none;background:#000;color:#fff;white-space:pre;font:10pt monospace'>")
			send(w, strings.Repeat("<!-- xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx -->\n", 13))
			send(w, "VERIFYING SECTOR 0x0000 <audio src='a/0x0000.mp3' autoplay volume='0.5'></audio>")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>USERNAME: ")
			send(w, "ZER0C00L")
			send(w, "<br>PASSWORD: ")
			send(w, "1")
			send(w, "5")
			send(w, "0")
			send(w, "7")
			<-time.After(time.Second * 3)
			send(w, "          [<font color='#0f0'>VALID</font>]<br>")
			send(w, "VERIFYING SECTOR 0x0001 <audio src='a/0x0001.mp3' autoplay volume='0.5'></audio>")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>ACCESS NUMBER:")
			send(w, "2")
			send(w, "1")
			send(w, "2")
			send(w, "5")
			send(w, "5")
			send(w, "5")
			send(w, "4")
			send(w, "2")
			send(w, "4")
			send(w, "0")
			<-time.After(time.Second * 18)
			send(w, "[<font color='#0f0'>VALID</font>]<br>")

			send(w, "VERIFYING SECTOR 0x0002 <audio src='a/0x0002.mp3' autoplay volume='0.5'></audio>")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>BLOCK: ")
			send(w, "A2")
			<-time.After(time.Second * 1)
			send(w, "               [<font color='#0f0'>PROMOTED TO ADVANCED ENGLISH</font>]<br>")

			t := uuid()
			tokens.Set(t, 0)
			send(w, `<meta http-equiv='refresh' content="2;URL='/?`+t+`=0'">`)
			return

		// SECTOR 1
		// SECTOR 1
		// SECTOR 1
		case "'sector 0x01'.torrent":
			head(w, 307, headers)
			send(w, "<html><body style='user-select:none;background:#000;color:#fff;white-space:pre;font:10pt monospace'><style>@keyframes blinkText{0%{color:#000}100%{color:#f00}}blink{animation:500ms blinkText infinite}</style>")
			send(w, strings.Repeat("<!-- xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx -->\n", 13))
			send(w, "VERIFYING SECTOR 0x0100 <audio src='a/0x0100.mp3' autoplay volume='0.5'></audio>")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>PLEASE DEPOSIT $5 ")
			send(w, ".")
			send(w, ".")
			send(w, ". ")
			send(w, "[<font color='#0f0'>THANK YOU</font>]<br>")

			send(w, "VERIFYING SECTOR 0x0101 <audio src='a/0x0101.mp3' autoplay volume='0.5'></audio>")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>CURRENT TIME: ")
			send(w, "0417")
			send(w, "    [<font color='#0f0'>VALID</font>]<br>")

			send(w, "VERIFYING SECTOR 0x0102 <audio src='a/0x0102.mp3' autoplay volume='0.5'></audio>")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>PASSPHRASE 1/3 ")
			send(w, "[<font color='#0f0'>VALID</font>]<br>")
			send(w, "<br>PASSPHRASE 2/3 ")
			send(w, "[<font color='#f00'><blink>MISSING</blink></font>]<br>")
			send(w, "<br>PASSPHRASE 3/3 ")
			send(w, "[<font color='#0f0'>VALID</font>]<br>")
			send(w, `<br><b>ENTER MISSING PASSPHRASE<blink>:</blink></b><br><form action='/0x01' method='post'><input type="text" name="passphrase" maxlength="45" placeholder='case insensitive/you insensitive twat'></form>`)
			/* continued above ... */
			return

		case "'sector 0x02'.torrent":
			head(w, 409, headers)
			send(w, "<html><body style='user-select:none;background:#000;color:#fff;white-space:pre;font:10pt monospace'><style>@keyframes passwordProtect{0%{background:#fff;color:#000;border-radius:0%}100%{background:#000;color:#000;border-radius:50%}}span{display:inline-block;animation:100ms passwordProtect;animation-delay:277ms;animation-fill-mode:forwards}@keyframes flash{0%{color:#000}100%{color:#0f0}}h1{animation:1s flash infinite}</style>")
			send(w, strings.Repeat("<!-- xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx -->\n", 13))
			send(w, "VERIFYING SECTOR 0x0200 <audio src='a/0x0200.mp3' autoplay volume='0.5'></audio>")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>PLEASE ENTER THE FOUR MOST COMMONLY USED PASSWORDS:<br>")
			send(w, "<span>L</span>")
			send(w, "<span>O</span>")
			send(w, "<span>V</span>")
			send(w, "<span>E</span>")
			send(w, "                                           [<font color='#0f0'>VALID</font>]<br>")
			<-time.After(time.Second * 3)
			send(w, "<span>S</span>")
			send(w, "<span>E</span>")
			send(w, "<span>X</span>")
			send(w, "                                            [<font color='#0f0'>VALID</font>]<br>")
			<-time.After(time.Second * 3)
			send(w, "<span>S</span>")
			send(w, "<span>3</span>")
			send(w, "<span>K</span>")
			send(w, "<span>R</span>")
			send(w, "<span>I</span>")
			send(w, "<span>T</span>")
			send(w, "                                         [<font color='#0f0'>VALID</font>]<br>")
			<-time.After(time.Second * 3)
			send(w, "<span>G</span>")
			send(w, "<span>O</span>")
			send(w, "<span>D</span>")
			send(w, "                                            [<font color='#0f0'>VALID</font>]<br>")
			<-time.After(time.Second)

			send(w, "VERIFYING SECTOR 0x0201 <audio src='a/0x0201.mp3' autoplay volume='0.5'></audio>")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, " PASSWORD: ")
			send(w, "<span>S</span>")
			send(w, "<span>H</span>")
			send(w, "<span>E</span>")
			send(w, "<span>I</span>")
			send(w, "<span>K</span>")
			send(w, "  [<font color='#0f0'>VALID</font>]<br>")
			<-time.After(time.Millisecond * 206)
			send(w, "VERIFYING SECTOR 0x0202 <audio src='a/0x0202.mp3' autoplay volume='0.5'></audio>")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, " PASSWORD: ")
			send(w, "<span>R</span>")
			send(w, "<span>E</span>")
			send(w, "<span>Z</span>")
			send(w, "<span>N</span>")
			send(w, "<span>O</span>")
			send(w, "<span>R</span>")
			send(w, " [<font color='#0f0'>VALID</font>]<br>")
			<-time.After(time.Second * 4)
			<-time.After(time.Millisecond * 80)
			send(w, "VERIFYING SECTOR 0x0203 <audio src='a/0x0203.mp3' autoplay volume='0.5' loop></audio>")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, " PASSWORD: ")
			send(w, "<span>A</span>")
			send(w, "<span>D</span>")
			send(w, "<span>D</span>")
			send(w, "<span>I</span>")
			send(w, "<span>C</span>")
			send(w, "<span>T</span>")
			send(w, " [<font color='#0f0'>VALID</font>]<br><br><h1 style='text-align:center;color:#0f0'>ACCESS GRANTED</h1>")

			t := uuid()
			tokens.Set(t, 2)
			send(w, `<meta http-equiv='refresh' content="3;URL='/?`+t+`=2'">`)
			return
		// SECTOR 3
		// SECTOR 3
		// SECTOR 3
		case "'sector 0x03'.torrent":
			head(w, 402, headers)
			send(w, "<html><body style='user-select:none;background:#000;color:#fff;font:10pt monospace'><style>@keyframes flash{0%{color:#fff}100%{color:#f00}}span{animation:75ms flash infinite}</style>")
			send(w, strings.Repeat("<!-- xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx -->\n", 13))
			if strings.Contains(r.Header.Get("User-Agent"), "Chrome/") {
				send(w, "<div style='white-space: normal; background:#fff; color: #222; font: 12pt sans-serif; padding: 1em; width: 50%; margin: 1em auto; border-radius: 25px; box-shadow: 12px 12px 0 rgba(255,255,255,0.7)'><b>It's official!</b> Fun is no longer allowed on the web! Bow to your corpo overlords and click this immersion breaking ugly piece of shit because autoplay doesn't work in Chrome unless you click somewhere on the page first. gj fuckin morons.<br><br>btw if you @ me with some condescending stackoverflow/freenode what are you trying to do bullshit I swear to Christmas I'll slap you with a trout.<br><br><b title='yeah yeah its too tl;dr you couldn't bother to give a toss blow me you careless tosser'>PRESS PLAY &darr;</b><br><audio src='a/0x0300.mp3' autoplay controls volume='0.5'></audio></div>")
			} else {
				send(w, "<audio src='a/0x0300.mp3' autoplay volume='0.5'></audio>")
			}
			send(w, "<br><div><a target='_blank' href='https://soundcloud.com/dayvonjersen/0x0300' title='Free (as in freemium) listen on sowcow'><img src='v/soundcloud-80x15.png'></a></div>")
			send(w, "VERIFYING SECTOR 0x0300")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "Power-ON Reset")
			<-time.After(time.Second * 7)
			send(w, `&emsp;&emsp;&emsp;&emsp;&emsp;SMCC SPARC-Station 10 UP/MP POST version 3.0 (8/15/93)<br><br><br><br>CPU_#0&emsp; TI, TMS390Z55(3.x)  1Mb External cache<br>
            <br>CPU_#1&emsp;******* NOT installed *******
            <br>CPU_#2&emsp;******* NOT installed *******
            <br>RAM_#_&emsp;******* NOT tripled   *******`)
			send(w, `<br><br>Allocating SRMMU Context Table
            <br>Setting SRMMU Context Register
            <br>Setting SRMMU Context Table Pointer Register
            <br>Allocating SRMMU Level 1 Table
            <br>Mapping RAM
            <br>Mapping ROM<br><br>ttya initialized`)
			send(w, `<br>SPARCstation 10 (1 X 390Z55)
            <br>ROM Rev. 2.14 272 MB memory installed, Serial #6C61696E.
            <br>Ethernet address 2B:69:73:6D:61:69, Host ID: 7761696675
            <br><br>The IDPROM contents are invalid`)
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<span>ARE YOU CHALLENGING ME?")
			send(w, "?")
			send(w, `‽</span><br><br>MESS WITH THE BEST DIE LIKE THE REST:<form action="/0x03" method="post"><input type="text" name="code" placeholder=">>>> 0VERRiDE C0DE <<<<"></form>`)
			// continued above
			return

		// SECTOR 4, 5, 6
		// SECTOR 4, 5, 6
		// SECTOR 4, 5, 6
		case "'sector 0x07'.torrent":
			head(w, 302, headers)
			send(w, "<html><body style='user-select:none;background:#000;color:#fff;white-space:pre;font:10pt monospace'><style>@keyframes flash{0%{color:#fff}100%{color:#f00}}span{animation:125ms flash infinite}</style>")
			send(w, strings.Repeat("<!-- xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx -->\n", 13))
			clip := ""
			switch rand.Intn(10) {
			case 0:
				clip = "0x0400"
			case 1:
				clip = "0x0401"
			case 2:
				clip = "0x0402"
			case 3:
				clip = "0x0500"
			case 4:
				clip = "0x0600"
			case 5:
				clip = "0x0601"
			case 6:
				clip = "0x0602"
			case 7:
				clip = "0x0603"
			case 8:
				clip = "0x0604"
			case 9:
				clip = "0x0605"
			}

			if strings.Contains(r.Header.Get("User-Agent"), "Chrome/") {
				send(w, "<div style='white-space: normal; background:#fff; color: #222; font: 12pt sans-serif; padding: 1em; width: 50%; margin: 1em auto; border-radius: 25px; box-shadow: 12px 12px 0 rgba(255,255,255,0.7)'><b>It's official!</b> Fun is no longer allowed on the web! Bow to your corpo overlords and click this immersion breaking ugly piece of shit because autoplay doesn't work in Chrome unless you click somewhere on the page first. gj fuckin morons.<br><br>btw if you @ me with some condescending stackoverflow/freenode what are you trying to do bullshit I swear to Christmas I'll slap you with a trout.<br><br><b title='yak yak yak, get a job'>&rarr; PRESS PLAY &rarr;</b><audio src='a/"+clip+".mp3' autoplay loop controls></audio></div>")
			} else {
				send(w, "<audio src='a/"+clip+".mp3' autoplay volume='0.5'></audio>")
			}

			send(w, "<br>VERIFYING SECTOR 0x0400 ")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>VERIFYING SECTOR 0x0401 ")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>VERIFYING SECTOR 0x0402 ")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>VERIFYING SECTOR 0x0500 ")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>VERIFYING SECTOR 0x0600 ")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>VERIFYING SECTOR 0x0601 ")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>VERIFYING SECTOR 0x0602 ")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>VERIFYING SECTOR 0x0603 ")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>VERIFYING SECTOR 0x0604 ")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>VERIFYING SECTOR 0x0605 ")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, ".")
			send(w, "<br>RABBIT IN THE ADMINISTRATION")
			send(w, "<br>FLU SHOT ADMINISTERED")
			send(w, "<br><span>KERNEL PANIC BRAIN CANCER DETECTED</span>")
			t := uuid()
			tokens.Set(t, 4)
			send(w, `<meta http-equiv='refresh' content="2;URL='/?`+t+`=4'">`)
			return

		// SECTOR 7
		// SECTOR 7
		// SECTOR 7
		case "'31337'.torrent":
			head(w, 410, headers)
			t := uuid()
			tokens.Set(t, 5)
			send(w, "<html><body style='user-select:none;background:#000;color:#fff;white-space:pre;font:10pt monospace'>")
			send(w, strings.Repeat("<!-- xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx -->\n", 13))
			send(w, "<br>VERiFYiNG SECT0R 0x0777 . . . <audio src='a/0x0777.mp3' autoplay volume='0.5' loop></audio>")
			send(w, "<br>&emsp;&emsp;YOU'RE <!-- a cheater cheater pumpkin eater --><a href='/?"+t+"=5' style='color:#fff!important;text-decoration:none!important;cursor:default!important'>ELITE</a> B)")
			return
		}
		head(w, 404, headers)
		send(w, "404 get fucked nerd")
		log.Printf("%v 404", r.URL)

	})

	listenAddr := fmt.Sprintf("%s:%d", addr, port)
	log.Println("listening on", listenAddr)
	// log.Fatalln(http.ListenAndServe(listenAddr, nil))
	log.Fatalln(http.ListenAndServeTLS(listenAddr, certFile, keyFile, nil))
}

func req(r *http.Request) {
	log.Println(fmt.Sprint(
		hostAddr(r), " <-> ", remoteAddr(r), "\n",
		r.Header.Get("User-Agent"), "\n",
		" -> ", r.Method, " ", r.URL, "\n",
	))
}

func hostAddr(r *http.Request) string {
	host := r.Header.Get("Host")
	if host != "" {
		return host
	}
	return r.Host
}

func remoteAddr(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		if strings.Contains(ip, ",") {
			return ip[:strings.LastIndex(ip, ",")]
		}
		return ip
	}
	return r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
}

func fileExists(filename string) bool {
	f, err := os.Open(filename)
	f.Close()
	if os.IsNotExist(err) {
		return false
	}
	checkErr(err)
	return true
}

func isDir(filename string) bool {
	finfo, err := os.Stat(filename)
	checkErr(err)
	return finfo.IsDir()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
