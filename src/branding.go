package main

import (
	"html"
	"html/template"
	"net/url"
	"os"
	"strings"
)

// Footer "hosted from" line and the browser-tab icon - both configurable so
// they can be flipped without code changes (set in .env, then
// `docker compose restart app`).
const (
	// HOSTED_FROM picks the footer line + logo:
	//   vps      (default) "Hosted from my VPS" + Debian logo
	//   homelab  "Hosted from my homelab" + Raspberry Pi logo
	//   none     no line at all
	EnvHostedFrom = "HOSTED_FROM"

	// FAVICON_EMOJI is the emoji used as the tab icon (default 👾). An emoji has
	// no background, so it blends into any tab colour. Set it to "none" to use
	// the plain /favicon.ico file instead.
	EnvFaviconEmoji     = "FAVICON_EMOJI"
	defaultFaviconEmoji = "👾"
	maxFaviconRunes     = 8 // an emoji with joiners/modifiers can be several code points
)

// Hosting is the footer's "hosted from" line.
type Hosting struct {
	Label   string // empty = show nothing
	Logo    string // versioned image URL
	LogoAlt string
}

// hostedFrom returns the footer line for the current HOSTED_FROM setting.
func hostedFrom() Hosting {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(EnvHostedFrom))) {
	case "homelab":
		return Hosting{Label: "Hosted from my homelab", Logo: logoPath, LogoAlt: "Raspberry Pi logo"}
	case "none":
		return Hosting{}
	}
	return Hosting{Label: "Hosted from my VPS", Logo: debianLogoPath, LogoAlt: "Debian logo"}
}

// faviconHref returns a data: URL for an SVG containing the configured emoji,
// or "" to fall back to /favicon.ico. The emoji is HTML-escaped and length
// limited, so a stray value in .env can't break the page.
func faviconHref() template.URL {
	emoji := strings.TrimSpace(os.Getenv(EnvFaviconEmoji))
	switch strings.ToLower(emoji) {
	case "":
		emoji = defaultFaviconEmoji
	case "none", "off":
		return ""
	}
	if r := []rune(emoji); len(r) > maxFaviconRunes {
		emoji = string(r[:maxFaviconRunes])
	}
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y=".9em" font-size="90">` +
		html.EscapeString(emoji) + `</text></svg>`
	// QueryEscape percent-encodes every reserved character (including & and ;);
	// it writes spaces as "+", which a data: URL would not decode, so use %20.
	encoded := strings.ReplaceAll(url.QueryEscape(svg), "+", "%20")
	return template.URL("data:image/svg+xml," + encoded)
}
