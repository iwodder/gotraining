package link

import (
	"errors"
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Link struct {
	Href string
	Text []string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, errors.New("invalid HTML was supplied")
	}
	links := make([]Link, 0)
	treeWalker(doc, &links)
	return links, nil
}

func treeWalker(n *html.Node, links *[]Link) {
	if n.Type == html.ElementNode && n.Data == "a" {
		l := Link{}
		for _, attrs := range n.Attr {
			k := strings.ToLower(attrs.Key)
			if k == "href" {
				l.Href = attrs.Val
			}
		}
		l.Text = extractText(n)
		*links = append(*links, l)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		treeWalker(c, links)
	}
}

func extractText(n *html.Node) []string {
	var result []string
	if n.Type == html.TextNode {
		trimmed := strings.TrimSpace(n.Data)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, extractText(c)...)
	}
	return result
}
