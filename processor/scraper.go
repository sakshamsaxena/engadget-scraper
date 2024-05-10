package processor

import (
	"net/http"
	"strings"
	"time"

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
		// TODO: retry better ?
		time.Sleep(time.Second)
		return scrapeURL(client, url)
	}

	mainNode, parseErr := html.Parse(response.Body)
	if parseErr != nil {
		panic(parseErr) // TODO: retry?
	}
	divNode := findNode(mainNode)
	essay := &strings.Builder{}
	extractText(divNode, essay)
	return essay.String()
}

func findNode(node *html.Node) *html.Node {
	if node.DataAtom == atom.Div {
		for _, a := range node.Attr {
			if a.Key == "class" && a.Val == "caas-body" {
				return node
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if divNode := findNode(child); divNode != nil {
			return divNode
		}
	}
	return nil
}

func extractText(node *html.Node, essay *strings.Builder) {
	if node == nil {
		return
	}
	if node.Type == html.TextNode {
		essay.WriteString(node.Data)
		return
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		extractText(child, essay)
	}
}
