package main

import (
	"flag"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	var src string
	var inc string
	var out string
	flag.StringVar(&src, "i", "src", "Input directory for pages to be built")
	flag.StringVar(&inc, "inc", "include", "Directory of imported templates")
	flag.StringVar(&out, "o", "out", "Output directory")
	flag.Parse()
	srcs := getPaths(src)
	tmpl, err := template.ParseFiles(append(srcs, getPaths(inc)...)...)
	if err != nil {
		panic(err)
	}
	os.RemoveAll(out)
	err = os.Mkdir(out, 0777)
	if err != nil {
		panic(err)
	}
	err = os.Mkdir(out+"/resources", 0777)
	if err != nil {
		panic(err)
	}
	tmpl.Funcs(map[string]interface{}{
		"Resource": func(src string) string {
			u, err := url.Parse(src)
			if err != nil {
				panic(err)
			}
			outpath := out + "/resources/" + filepath.Base(u.Path)
			if _, err = os.Stat(outpath); err == nil {
				return outpath[len(out):]
			}
			o, err := os.OpenFile(outpath, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				panic(err)
			}
			defer func() {
				e := o.Close()
				if e != nil {
					panic(e)
				}
			}()
			g, err := http.Get(u.String())
			if err != nil {
				panic(err)
			}
			defer func() {
				e := g.Body.Close()
				if e != nil {
					panic(e)
				}
			}()
			_, err = io.Copy(o, g.Body)
			if err != nil {
				panic(err)
			}
			return outpath[len(out):]
		},
	})
	for _, v := range srcs {
		func() {
			os.MkdirAll(filepath.Dir(out+"/"+v[len(src)+1:]), 0777)
			ofile, err := os.OpenFile(out+"/"+v[len(src)+1:], os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				panic(err)
			}
			defer func() {
				e := ofile.Close()
				if e != nil {
					panic(e)
				}
			}()
			err = tmpl.ExecuteTemplate(ofile, filepath.Base(v[len(src)+1:]), map[string]interface{}{"Path": v[len(src)+1:]})
			if err != nil {
				panic(err)
			}
		}()
	}
}

func getPaths(dir string) (o []string) {
	o = []string{}
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	if info.IsDir() {
		infos, _ := f.Readdir(0)
		for _, v := range infos {
			o = append(o, getPaths(dir+"/"+v.Name())...)
		}
	} else {
		o = append(o, dir)
	}
	return
}
