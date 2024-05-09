package workers

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func scrapeURL(client *http.Client, url string) string {
	response, opErr := client.Get(url)
	if opErr != nil {
		panic(opErr) // TODO: retry
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		// TODO: retry
	}

	mainNode, parseErr := html.Parse(response.Body)
	if parseErr != nil {
		panic(parseErr) // TODO: retry
	}
	divNode := f(mainNode)
	essay := &strings.Builder{}
	t(divNode, essay)
	//fmt.Println(url)
	//fmt.Println(essay.String())
	return essay.String()
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
