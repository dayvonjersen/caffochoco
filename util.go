package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	pp "github.com/maruel/panicparse/stack"
)

func checkErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
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

func getStack(stack []byte) string {
	in := bytes.NewBuffer(stack)
	trace, err := pp.ParseDump(in, ioutil.Discard)
	checkErr(err)
	p := &pp.Palette{}
	buckets := pp.SortBuckets(pp.Bucketize(trace, pp.AnyValue))
	src, pkg := pp.CalcLengths(buckets, false)
	ret := ""
	for _, bucket := range buckets {
		ret += p.StackLines(&bucket.Signature, src, pkg, false)
	}
	return ret
}

func renderTemplate(filename string, data interface{}) string {
	t, err := template.ParseFiles(filename)
	checkErr(err)
	buf := new(bytes.Buffer)
	checkErr(t.Execute(buf, data))
	return buf.String()
}
