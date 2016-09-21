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

var imageFSHandler = fasthttp.FSHandler(IMAGE_DIR, 1)

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
			fmt.Fprintln(ctx, swatch.String())
		}
	} else {
		imageFSHandler(ctx)
	}

}
