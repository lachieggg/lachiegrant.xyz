package xml

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const titleOpen = "<title>"
const titleClose = "</title>"
const titleString = titleOpen + "%s" + titleClose

// MergeXML
func MergeXML(firstXML []byte, secondXML []byte, title string) string {
	// Extract contents inside <body> tags for both XML outputs
	extractBodyContent := func(xmlData []byte) string {
		start := bytes.Index(xmlData, []byte("<body>"))
		end := bytes.Index(xmlData, []byte("</body>"))
		if start == -1 || end == -1 {
			return ""
		}
		return string(xmlData[start+len("<body>") : end])
	}

	firstBody := extractBodyContent(firstXML)
	secondBody := extractBodyContent(secondXML)

	// Combine the two extracted body contents
	combined := firstBody + "<br><br>" + secondBody

	header := "<head><title>Status</title></head>"
	body := "<body>%s</body>"
	// Create the final merged XML
	mergedXML := fmt.Sprintf("<html>"+header+body+"</html>", combined)

	return mergedXML
}

// Replacer
func Replacer(input string) string {
	return strings.Replace(
		input,
		fmt.Sprintf(titleString, "stdin"),
		fmt.Sprintf(titleString, "Status"),
		1,
	)
}

// stripNewlinesAndTabs
func stripNewlinesAndTabs(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

// extractBodyContent extracts content inside the <body> tags from a given HTML string.
func extractBodyContent(htmlStr string) (content string, err error) {
	re := regexp.MustCompile(`(?s)<body.*?>(.*?)<\/body>`)
	matches := re.FindStringSubmatch(htmlStr)
	if len(matches) < 2 {
		return "", errors.New("could not find content within <body> tags")
	}
	return matches[1], nil
}

// MergeBodyContents takes in two HTML strings and merges their body contents.
func MergeBodyContents(html1, html2 string) (merged string, err error) {
	content1, err1 := extractBodyContent(html1)
	content2, err2 := extractBodyContent(html2)
	if err1 != nil || err2 != nil {
		return "", errors.New("error extracting body contents")
	}

	header := "<html><head><title>Status</title></head>"
	return header + fmt.Sprintf("<body>%s<br><br>%s</body>", content1, content2) + "</html>", nil
}
