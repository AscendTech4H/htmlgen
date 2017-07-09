package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
		titlecp := make([]interface{}, strings.Count(tmpl, "%s")-1)
		for i := range titlecp {
			titlecp[i] = title
		}
		pg := fmt.Sprintf(tmpl, append(titlecp, string(dat))...)
		pg = strings.Replace(pg, "\n", "", -1)
		pg = strings.Replace(pg, "\t", "", -1)
		chk(ioutil.WriteFile(opath, []byte(pg), os.ModePerm))
	}
}
