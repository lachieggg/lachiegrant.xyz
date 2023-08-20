package xml

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

const tagErrorMsg = "Error extracting tags: %v"

// ExtractTags
func ExtractTags(htmlContent string, tagName string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var tags []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tagName {
			var b strings.Builder
			html.Render(&b, n)
			tags = append(tags, b.String())
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return tags, nil
}

// MergeHTMLContents
func MergeHTMLContents(htmlContent1, htmlContent2 string) (string, error) {
	bodies1, err := ExtractTags(htmlContent1, "body")
	if err != nil || len(bodies1) == 0 {
		return "", fmt.Errorf(tagErrorMsg, err)
	}

	bodies2, err := ExtractTags(htmlContent2, "body")
	if err != nil || len(bodies2) == 0 {
		return "", fmt.Errorf(tagErrorMsg, err)
	}

	styles1, err := ExtractTags(htmlContent1, "style")
	if err != nil {
		return "", fmt.Errorf(tagErrorMsg, err)
	}

	styles2, err := ExtractTags(htmlContent2, "style")
	if err != nil {
		return "", fmt.Errorf(tagErrorMsg, err)
	}

	return fmt.Sprintf("<html><head>%s%s</head>%s%s</html>",
		strings.Join(styles1, ""),
		strings.Join(styles2, ""),
		bodies1[0],
		bodies2[0],
	), nil
}
