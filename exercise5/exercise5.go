package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"gotraining/exercise4/link"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc string `xml:"loc"`
}

func main() {
	urlArg := flag.String("url", "https://gophercises.com", "The URL which needs a sitemap generated.")
	maxDepth := flag.Int("depth", 5, "How far to search from the first page.")
	flag.Parse()

	urls := searchSite(*urlArg, *maxDepth)

	urlSet := URLSet{
		URLs: []URL{},
	}

	for _, v := range urls {
		urlSet.URLs = append(urlSet.URLs, URL{
			Loc: v,
		})
	}

	bytes, err := xml.MarshalIndent(&urlSet, "", "  ")
	if err != nil {
		log.Fatal("Uh-oh")
	}
	fmt.Println(string(bytes))
}

func searchSite(url string, maxDepth int) []string {
	queue := make(map[string]struct{})
	nextQueue := map[string]struct{}{
		url: {},
	}
	visitedPages := make(map[string]struct{})

	for i := 0; i <= maxDepth; i++ {
		queue, nextQueue = nextQueue, make(map[string]struct{})
		for page, _ := range queue {
			visitedPages[page] = struct{}{}
			for _, childPage := range getPageLinks(page) {
				if _, ok := visitedPages[childPage]; ok {
					continue
				}
				nextQueue[childPage] = struct{}{}
			}
		}
	}

	ret := make([]string, 0, len(visitedPages))
	for k, _ := range visitedPages {
		ret = append(ret, k)
	}
	return ret
}

func getPageLinks(page string) []string {
	resp, err := http.Get(page)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		return nil
	}
	if err != nil {
		return nil
	}
	baseUrl := url.URL{
		Scheme: resp.Request.URL.Scheme,
		Host:   resp.Request.URL.Host,
	}
	return extractLinks(resp.Body, baseUrl.String())
}

func extractLinks(r io.Reader, baseUrl string) []string {
	var ret []string
	links, err := link.Parse(r)
	if err != nil {
		return ret
	}
	for _, l := range links {
		if strings.HasPrefix(l.Href, "/") {
			l.Href = baseUrl + l.Href
		}
		if strings.HasPrefix(l.Href, baseUrl) {
			ret = append(ret, l.Href)
		}
	}
	return ret
}

// 1. Get base URL from user
// 2. Make GET request to the base URL
// 3. Parse returned HTML for links
// 4. Add returned links to some structure.... tree?
