package main

import (
	"fmt"
	"github.com/jackdanger/collectlinks"
	"net/http"
	"net/url"
)

var visited = make(map[string]bool)

func main() {
	fmt.Println("hello world")
	url := "https://blog.onepie.cn"
	queue := make(chan string)
	go func() {
		queue <- url
	}()
	for uri := range queue {
		getPageUrls(uri, queue)
	}
}

// 解析网页
// 常见解析器xpath jquery
// 直接使用轮子 github.com/jackdanger/collectlinks
func getPageUrls(url string, queue chan string) {
	visited[url] = true
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	//自定义header
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return
	}
	defer resp.Body.Close()
	links := collectlinks.All(resp.Body)

	for _, link := range links {
		absolute := urlJoin(link, url)
		if url != " " {
			if !visited[absolute] {
				fmt.Println("paese url", absolute)
				go func() {
					queue <- absolute
				}()
			}
		}
	}
}

// 解析url 拼接完整地址
func urlJoin(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	return baseUrl.ResolveReference(uri).String()
}
