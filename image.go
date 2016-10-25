package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"github.com/generaltso/vibrant"
	"github.com/valyala/fasthttp"
)

var imageFSHandler = fasthttp.FSHandler("./image", 1)

func imageHandler(ctx *fasthttp.RequestCtx) {
	path := "." + string(ctx.Path())
	if strings.HasSuffix(path, ".css") {
		path = strings.TrimSuffix(path, ".css")
		if !fileExists(path) {
			notfoundHandler(ctx)
			return
		}
		f, err := os.Open(path)
		checkErr(err)
		img, _, err := image.Decode(f)
		f.Close()
		checkErr(err)
		palette, err := vibrant.NewPaletteFromImage(img)
		checkErr(err)
		ctx.Response.Header.Set("Content-Type", "text/css; charset=UTF-8")
		for _, swatch := range palette.ExtractAwesome() {
			c := swatch.Color
			r, g, b := c.RGB()
			fmt.Fprintf(ctx, `
.%s {
	background-color: %s;
	color: %s;
	box-shadow: 0 0 10px rgba(%d,%d,%d,1); 
}`, strings.ToLower(swatch.Name), c, c.TitleTextColor(), r, g, b)
		}
	} else if path == "./image/" || !fileExists(path) {
		notfoundHandler(ctx)
	} else {
		imageFSHandler(ctx)
	}

}
