package restaurants

import (
	"errors"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
)

func getAttribute(node *html.Node, key string) (string, error) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val, nil
		}
	}

	return "", errors.New("couldn't find the provided key")
}

func hasClass(node *html.Node, class string) bool {
	if node.Type == html.ElementNode {
		c, err := getAttribute(node, "class")
		if err != nil {
			return false
		}
		return c == class
	}
	return false
}

func findNodeByClass(node *html.Node, class string) (*html.Node, error) {
	if hasClass(node, class) {
		return node, nil
	}
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		c, err := findNodeByClass(n, class)
		if err == nil {
			return c, nil
		}
	}
	return nil, errors.New("couldn't find a node with provided class")
}

func getTextInternal(node *html.Node) (string, error) {
	if node.Type == html.TextNode {
		return node.Data, nil
	}
	return "", errors.New("not a text node")
}

func getText(node *html.Node) (string, error) {
	if node.Type == html.TextNode {
		return node.Data, nil
	}
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		text, err := getTextInternal(n)
		if err == nil {
			return text, nil
		}
	}
	return "", errors.New("couldn't find a text node")
}

func getTextDecodeWindows1250(node *html.Node) (string, error) {
	text, err := getText(node)
	if err != nil {
		return text, err
	}
	return decodeWindows1250(text)
}

func decodeWindows1250(text string) (string, error) {
	dec := charmap.Windows1250.NewDecoder()
	out, err := dec.String(text)
	return out, err
}
