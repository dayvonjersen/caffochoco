package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strings"

	"github.com/generaltso/vibrant"
)

var imageFSHandler = http.StripPrefix("/image/", http.FileServer(http.Dir("./image")))

func imageHandler(w http.ResponseWriter, r *http.Request) {
	path := "." + r.URL.Path
	if strings.HasSuffix(path, ".css") {
		path = strings.TrimSuffix(path, ".css")
		if !fileExists(path) {
			notfoundHandler(w, r)
			return
		}
		f, err := os.Open(path)
		checkErr(err)
		img, _, err := image.Decode(f)
		f.Close()
		checkErr(err)
		palette, err := vibrant.NewPaletteFromImage(img)
		checkErr(err)
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		for _, swatch := range palette.ExtractAwesome() {
			c := swatch.Color
			r, g, b := c.RGB()
			fmt.Fprintf(w, `
.%s {
	background-color: %s;
	color: %s;
/*	box-shadow: 0 0 10px rgba(%d,%d,%d,1); */
}`, strings.ToLower(swatch.Name), c, c.TitleTextColor(), r, g, b)
		}
	} else if path == "./image/" || !fileExists(path) {
		notfoundHandler(w, r)
	} else {
		imageFSHandler.ServeHTTP(w, r)
	}

}
