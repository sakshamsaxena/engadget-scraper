package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Scraper struct {
	client *http.Client
}

func New(client *http.Client) *Scraper {
	return &Scraper{
		client: client,
	}
}

func (s *Scraper) Scrape(url string) {
	resp, err := s.client.Get(url)
	if err != nil {
		panic(err) // TODO: handle
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	n, e := html.Parse(strings.NewReader(string(body)))
	if e != nil {
		panic(e)
	}
	divNode := f(n)
	essay := &strings.Builder{}
	t(divNode, essay)
	fmt.Println(url)
	fmt.Println(essay.String())
}

func f(n *html.Node) *html.Node {
	if n.DataAtom == atom.Div {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "caas-body" {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if k := f(c); k != nil {
			return k
		}
	}
	return nil
}

func t(n *html.Node, essay *strings.Builder) {
	if n == nil {
		fmt.Println("nil node")
		return
	}
	if n.Type == html.TextNode {
		essay.WriteString(n.Data)
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		t(c, essay)
	}
}
