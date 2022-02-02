package sitemapbuilder

import (
	"encoding/xml"
	htmlparser "github.com/gmaschi/go-sitemap-builder/html-parser"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

type (
	loc struct {
		Value string `xml:"loc"`
	}

	urlset struct {
		Urls  []loc  `xml:"url"`
		Xmlns string `xml:"xmlns,attr"`
	}
)

func Build(urlPath string, maxDepth int) ([]byte, error) {
	pages := bfs(urlPath, maxDepth)
	data, err := toXml(pages)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func toXml(pages []string) ([]byte, error) {
	xmlData := urlset{
		Urls:  make([]loc, 0, len(pages)),
		Xmlns: xmlns,
	}

	for _, page := range pages {
		xmlData.Urls = append(xmlData.Urls, loc{Value: page})
	}

	data, err := xml.MarshalIndent(xmlData, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), data...), nil
}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}

	for i := 0; i < maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			links, err := get(url)
			if err != nil {
				continue
			}
			for _, l := range links {
				if _, ok := seen[l]; !ok {
					nq[l] = struct{}{}
				}
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for i := range seen {
		ret = append(ret, i)
	}
	return ret
}

func get(urlStr string) ([]string, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	return hrefs(resp.Body, base)
}

func hrefs(r io.Reader, base string) ([]string, error) {
	links, _ := htmlparser.Parse(r)
	ret := make([]string, 0, len(links)/2)
	for _, l := range links {
		switch {
		case strings.Contains(l.Href, "#") ||
			strings.Contains(l.Href, "?"):
			continue
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http") &&
			strings.Contains(l.Href, base):
			ret = append(ret, l.Href)
		}
	}
	return ret, nil
}
