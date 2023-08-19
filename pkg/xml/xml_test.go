package xml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMergeXML
func TestMergeXML(t *testing.T) {
	htopMock := []byte("<html><head></head><body>htop</body></html>")
	neofetchMock := []byte("<html><head></head><body>neofetch</body></html>")

	merged := MergeXML(htopMock, neofetchMock, "Status")
	assert.Contains(t, merged, "htop")
	assert.Contains(t, merged, "neofetch")

	header := "<head><title>Status</title></head>"
	body := "<body>htop<br><br>neofetch</body>"
	expected := "<html>" + header + body + "</html>"

	expected = stripNewlinesAndTabs(expected)
	merged = stripNewlinesAndTabs(merged)

	assert.Equal(t, expected, merged)
}
