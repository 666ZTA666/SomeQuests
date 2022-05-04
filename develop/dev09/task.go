package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	R := flag.Bool("r", false, "рекурсия")
	V := flag.Bool("v", false, "для скачивания css")

	flag.Parse()
	Wurl := flag.Arg(0)

	resp := checkGet(Wurl)
	if resp == nil {
		fmt.Println("cant make GET to url", Wurl)
		os.Exit(0)
	}
	r, f := outInFile(Wurl, resp)

	if *R {
		recursiveShit(Wurl, r, *V)
	} else {
		data, err := io.ReadAll(r)
		if err != nil {
			fmt.Println("reading from r err", err)
		}
		ss := strings.Split(string(data), "\n")
		for _, s := range ss {
			fmt.Println(s)
		}
	}

	err := resp.Body.Close()
	if err != nil {
		fmt.Println("resp.body closing err", err)
	}
	err = f.Close()
	if err != nil {
		fmt.Println("close file err", err)
	}

}

func outInFile(hostPath string, resp *http.Response) (io.Reader, *os.File) {
	//todo сделать подпись для возможности перехода по ссылкам локально
	var filename string
	if hostPath == resp.Request.URL.Host {
		err := os.Mkdir(resp.Request.URL.Host, os.FileMode(777))
		if err != nil {
			fmt.Println("one dir creating err")
			fmt.Println(err)
		}
		filename = "index.html"
	} else {
		path := resp.Request.URL.Path
		if len(path) == 1 {
			return nil, nil
		}

		if path[len(path)-1] == '/' {
			filename = path[1:] + "index.html"
		} else {
			filename = path[1:]
		}
	}
	stat, err := os.Stat(filepath.Dir(resp.Request.URL.Host + "/" + filename))
	if err == nil && !stat.IsDir() {
		err = os.Remove(filepath.Dir(resp.Request.URL.Host + "/" + filename))
		if err != nil {
			fmt.Println("dir remove err", err)
		}
	}
	err = os.MkdirAll(filepath.Dir(resp.Request.URL.Host+"/"+filename), os.FileMode(666))
	if err != nil {
		fmt.Println("all dir creating error", err)
		fmt.Println(filepath.Dir(resp.Request.URL.Host + "/" + filename))
	}
	file, err := os.Create(resp.Request.URL.Host + "/" + filename)
	if err != nil {
		fmt.Println("file creating error", err)
		fmt.Println(resp.Request.URL.Host, "/", filename)
	}
	r := io.TeeReader(resp.Body, file)
	return r, file
}

func checkGet(s string) *http.Response {
	resp, err := http.Get("http://" + s)
	if err != nil {
		fmt.Println("GET to this url err", err)
		return nil
	}
	if resp.StatusCode != 200 {
		fmt.Println("status code != 200")
		return nil
	}
	return resp
}

func recursiveShit(s string, ir io.Reader, v bool) {
	zz := getLink(ir, v, s)

	for _, link := range zz {
		resp := checkGet(s + link)
		if resp == nil {
			fmt.Println("cant make GET to url", s+link)
			continue
		}
		or, f := outInFile(s+link, resp)

		recursiveShit(s+link, or, v)

		err := resp.Body.Close()
		if err != nil {
			fmt.Println("resp.body closing err", err)
		}
		err = f.Close()
		if err != nil {
			fmt.Println("close file err", err)
		}
	}
}

func getLink(body io.Reader, v bool, hostPath string) []string {
	//todo сделать скачивание скриптов
	hash := make(map[string]struct{})
	t := html.NewTokenizer(body)
	for {
		tn := t.Next()
		if tn == html.ErrorToken {
			break
		} else if tn == html.StartTagToken || tn == html.EndTagToken {
			token := t.Token()
			if token.Data == "link" && v {
				for _, attr := range token.Attr {
					if attr.Key == "href" && strings.Contains(attr.Val, ".css") {
						urlStr, err := url.Parse(attr.Val)
						if err != nil {
							continue
						}
						if urlStr.Host == "" {
							hash[strings.TrimSpace(urlStr.Path)] = struct{}{}
						} else if !strings.Contains(urlStr.Host+urlStr.Path, hostPath) {
							continue
						}
						hash[strings.TrimSpace(strings.TrimLeft(urlStr.Host+urlStr.Path, hostPath))] = struct{}{}
					}
				}
			} else if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						urlStr, err := url.Parse(attr.Val)
						if err != nil {
							continue
						}
						if urlStr.Host == "" {
							hash[strings.TrimSpace(urlStr.Path)] = struct{}{}
						} else if !strings.Contains(urlStr.Host+urlStr.Path, hostPath) {
							continue
						}
						hash[strings.TrimSpace(strings.TrimLeft(urlStr.Host+urlStr.Path, hostPath))] = struct{}{}
					}
				}
			}
		}
	}
	links := make([]string, 0, len(hash))
	i := 0
	for k := range hash {
		if k != "" {
			links = append(links, k)
			i++
		}
	}
	return links
}
