package xml

import (
	"bytes"
	"fmt"
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

	htopBody := extractBodyContent(firstXML)
	neofetchBody := extractBodyContent(secondXML)

	// Combine the two extracted body contents
	combined := htopBody + "<br><br>" + neofetchBody

	// Create the final merged XML
	mergedXML := fmt.Sprintf(`
		<html>
		<head>
			<title>%s</title>
		</head>
		<body>
		%s
		</body>
		</html>
	`, title, combined)

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
