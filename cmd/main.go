package main

import (
	"flag"
	sitemapbuilder "github.com/gmaschi/go-sitemap-builder/sitemap-builder"
	"os"
)

func main() {
	url := flag.String("url", "https://www.psymeetsocial.com", "url to build the sitemap from")
	maxDepth := flag.Int("depth", 20, "max depth to traverse site links")
	flag.Parse()

	data, err := sitemapbuilder.Build(*url, *maxDepth)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("sitemap.xml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		panic(err)
	}
}
