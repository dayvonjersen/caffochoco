package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/exp/inotify"
)

type task struct {
	cmd, dir string
}

func (t *task) run() error {
	log.Println("executing", t.cmd, "from dir", t.dir)

	cmd := exec.Command("sh", "-c", t.cmd)
	cmd.Dir = t.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	start := time.Now()
	err := cmd.Run()
	end := time.Now()

	log.Println("finished", t.cmd, end.Sub(start))
	return err
}

func main() {
	w, err := inotify.NewWatcher()
	checkErr(err)

	const appDir = "./app/"
	checkErr(w.Watch(appDir))
	filepath.Walk(appDir, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			checkErr(w.Watch(path))
		}
		return nil
	})

	queue := map[string]struct{}{}
	tasks := make(chan *task)
	go func() {
		for {
			t := <-tasks
			err := t.run()
			if err != nil {
				log.Println(err)
			}
			delete(queue, t.dir)
		}
	}()

	for {
	here:
		select {
		case e := <-w.Event:
			evtdir := strings.Replace(filepath.Dir(e.Name), "\\", "/", -1)
			rundir := strings.Split(evtdir, "/")[0]
			file := filepath.Base(e.Name)
			ext := filepath.Ext(file)
			cmd := "./build.sh"

			if len(file) == 0 || len(ext) == 0 || file[0:1] == "." || ext[len(ext)-1:] == "~" {
				break here
			}
			switch ext {
			case ".swp":
				fallthrough
			case ".swo":
				fallthrough
			case ".swn":
				fallthrough
			case ".tmp":
				break here
			}

			if _, ok := queue[rundir]; !ok {
				queue[rundir] = struct{}{}
				tasks <- &task{cmd, rundir}
			}
		case err := <-w.Error:
			checkErr(err)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
