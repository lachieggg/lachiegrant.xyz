package xml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractTags_SingleTag(t *testing.T) {
	html := `<html><body><div>Hello</div></body></html>`

	tags, err := ExtractTags(html, "div")

	require.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, "<div>Hello</div>", tags[0])
}

func TestExtractTags_MultipleTags(t *testing.T) {
	html := `<html><body><p>First</p><p>Second</p><p>Third</p></body></html>`

	tags, err := ExtractTags(html, "p")

	require.NoError(t, err)
	assert.Len(t, tags, 3)
	assert.Equal(t, "<p>First</p>", tags[0])
	assert.Equal(t, "<p>Second</p>", tags[1])
	assert.Equal(t, "<p>Third</p>", tags[2])
}

func TestExtractTags_NestedTags(t *testing.T) {
	html := `<html><body><div><span>Nested</span></div></body></html>`

	// Should find the outer div (which contains the span)
	divTags, err := ExtractTags(html, "div")
	require.NoError(t, err)
	assert.Len(t, divTags, 1)
	assert.Contains(t, divTags[0], "<span>Nested</span>")

	// Should also find the nested span
	spanTags, err := ExtractTags(html, "span")
	require.NoError(t, err)
	assert.Len(t, spanTags, 1)
	assert.Equal(t, "<span>Nested</span>", spanTags[0])
}

func TestExtractTags_NoMatchingTags(t *testing.T) {
	html := `<html><body><div>Content</div></body></html>`

	tags, err := ExtractTags(html, "span")

	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestExtractTags_BodyTag(t *testing.T) {
	html := `<html><body><p>Content</p></body></html>`

	tags, err := ExtractTags(html, "body")

	require.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Contains(t, tags[0], "<p>Content</p>")
}

func TestExtractTags_StyleTag(t *testing.T) {
	html := `<html><head><style>.red { color: red; }</style></head><body></body></html>`

	tags, err := ExtractTags(html, "style")

	require.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Contains(t, tags[0], ".red { color: red; }")
}

func TestExtractTags_InvalidHTML(t *testing.T) {
	// html.Parse is very forgiving, so even malformed HTML typically parses
	html := `<div>unclosed`

	tags, err := ExtractTags(html, "div")

	require.NoError(t, err)
	assert.Len(t, tags, 1)
}

func TestExtractTags_EmptyInput(t *testing.T) {
	tags, err := ExtractTags("", "div")

	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestMergeHTMLContents_Success(t *testing.T) {
	html1 := `<html><head><style>.a{}</style></head><body><p>Body1</p></body></html>`
	html2 := `<html><head><style>.b{}</style></head><body><p>Body2</p></body></html>`

	merged, err := MergeHTMLContents(html1, html2)

	require.NoError(t, err)
	assert.Contains(t, merged, "<html>")
	assert.Contains(t, merged, "</html>")
	assert.Contains(t, merged, ".a{}")
	assert.Contains(t, merged, ".b{}")
	assert.Contains(t, merged, "Body1")
	assert.Contains(t, merged, "Body2")
}

func TestMergeHTMLContents_NoStyleTags(t *testing.T) {
	html1 := `<html><body><p>Body1</p></body></html>`
	html2 := `<html><body><p>Body2</p></body></html>`

	merged, err := MergeHTMLContents(html1, html2)

	require.NoError(t, err)
	assert.Contains(t, merged, "Body1")
	assert.Contains(t, merged, "Body2")
}

func TestMergeHTMLContents_PreservesBodyOrder(t *testing.T) {
	// The first body should appear before the second body in the output
	html1 := `<html><body><div id="first">First</div></body></html>`
	html2 := `<html><body><div id="second">Second</div></body></html>`

	merged, err := MergeHTMLContents(html1, html2)

	require.NoError(t, err)
	// First body content should appear before second
	firstIdx := strings.Index(merged, "First")
	secondIdx := strings.Index(merged, "Second")
	assert.True(t, firstIdx < secondIdx, "first body should appear before second body")
}

func TestMergeHTMLContents_MultipleStyles(t *testing.T) {
	html1 := `<html><head><style>.a{}</style><style>.b{}</style></head><body><div></div></body></html>`
	html2 := `<html><head><style>.c{}</style></head><body><div></div></body></html>`

	merged, err := MergeHTMLContents(html1, html2)

	require.NoError(t, err)
	assert.Contains(t, merged, ".a{}")
	assert.Contains(t, merged, ".b{}")
	assert.Contains(t, merged, ".c{}")
}
