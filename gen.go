package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"strings"
	"os"
)

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var serve string
	var template string
	var out string
	flag.StringVar(&serve, "serve", "", "Serve HTTP at address")
	flag.StringVar(&template, "t", "", "Template to use for pages")
	flag.StringVar(&out, "o", "out", "Destination for files")
	flag.Parse()
	if template == "" {
		fmt.Println("Missing template")
		flag.PrintDefaults()
		return
	}
	files := []string{}
	args := flag.Args()
	for _, v := range args {
		matches, err := filepath.Glob(v)
		chk(err)
		files = append(files, matches...)
	}
	td, err := ioutil.ReadFile(template)
	chk(err)
	tmpl := string(td)
	for _, v := range files {
		dat, err := ioutil.ReadFile(v)
		chk(err)
		opath := out + "/" + v
		title := strings.Split(filepath.Base(v), ".")[0]
		chk(ioutil.WriteFile(opath, []byte(fmt.Sprintf(tmpl, title, string(dat))), os.ModePerm))
	}
}
