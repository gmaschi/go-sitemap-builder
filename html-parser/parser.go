package htmlparser

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

// Link represents a link (<a href="...">) in an HTML file
type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	//var links []Link
	//var i int

	//var fnText func(node *html.Node, i int)
	//fnText = func(node *html.Node, i int) {
	//	if node.Type == html.TextNode {
	//		links[i].Text += node.Data
	//	}
	//	for c := node.FirstChild; c != nil; c = c.NextSibling {
	//		fnText(c, i)
	//	}
	//}
	//
	//var fn func(node *html.Node)
	//fn = func(node *html.Node) {
	//	if node.Type == html.ElementNode && node.Data == "a" {
	//		for _, attr := range node.Attr {
	//			if attr.Key == "href" {
	//				links = append(links, Link{Href: attr.Val})
	//				fnText(node, i)
	//				i++
	//			}
	//		}
	//	}
	//	for c := node.FirstChild; c != nil; c = c.NextSibling {
	//		fn(c)
	//	}
	//}
	//fn(doc)

	nodes := nodeLinks(doc)
	links := make([]Link, 0, len(nodes))
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func nodeLinks(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}
	var ret []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, nodeLinks(c)...)
	}
	return ret
}

func buildLink(node *html.Node) Link {
	var ret Link
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = getNodeText(node)
	return ret
}

func getNodeText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	var ret string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret += getNodeText(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}
