package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func renderPage(t *testing.T, serve http.HandlerFunc, path string) string {
	t.Helper()
	rec := httptest.NewRecorder()
	serve(rec, httptest.NewRequest(http.MethodGet, path, nil))
	assert.Equal(t, http.StatusOK, rec.Code, path)
	return rec.Body.String()
}

func TestHostedFrom(t *testing.T) {
	t.Run("default is the VPS with the Debian logo", func(t *testing.T) {
		t.Setenv(EnvHostedFrom, "")
		h := hostedFrom()
		assert.Equal(t, "Hosted from my VPS", h.Label)
		assert.Equal(t, debianLogoPath, h.Logo)
		assert.Equal(t, "Debian logo", h.LogoAlt)
	})
	t.Run("homelab brings back the Raspberry Pi", func(t *testing.T) {
		for _, v := range []string{"homelab", " HomeLab "} {
			t.Setenv(EnvHostedFrom, v)
			h := hostedFrom()
			assert.Equal(t, "Hosted from my homelab", h.Label, v)
			assert.Equal(t, logoPath, h.Logo, v)
		}
	})
	t.Run("none hides the line, and an unknown value falls back to the default", func(t *testing.T) {
		t.Setenv(EnvHostedFrom, "none")
		assert.Empty(t, hostedFrom().Label)
		t.Setenv(EnvHostedFrom, "spaceship")
		assert.Equal(t, "Hosted from my VPS", hostedFrom().Label)
	})
}

func TestFooterReflectsHostedFrom(t *testing.T) {
	t.Run("vps", func(t *testing.T) {
		t.Setenv(EnvHostedFrom, "")
		body := renderPage(t, codeHandler, "/code")
		assert.Contains(t, body, "Hosted from my VPS")
		assert.Contains(t, body, `alt="Debian logo"`)
		assert.NotContains(t, body, "homelab")
	})
	t.Run("homelab", func(t *testing.T) {
		t.Setenv(EnvHostedFrom, "homelab")
		body := renderPage(t, codeHandler, "/code")
		assert.Contains(t, body, "Hosted from my homelab")
		assert.Contains(t, body, `alt="Raspberry Pi logo"`)
	})
	t.Run("none", func(t *testing.T) {
		t.Setenv(EnvHostedFrom, "none")
		body := renderPage(t, codeHandler, "/code")
		assert.NotContains(t, body, "Hosted from")
		assert.NotContains(t, body, `id="logo-img"`)
	})
}

func TestFaviconHref(t *testing.T) {
	t.Run("default is the house emoji as an SVG data URL", func(t *testing.T) {
		t.Setenv(EnvFaviconEmoji, "")
		href := string(faviconHref())
		assert.True(t, strings.HasPrefix(href, "data:image/svg+xml,"))
		assert.Contains(t, href, "%F0%9F%91%BE") // 👾
	})
	t.Run("custom emoji", func(t *testing.T) {
		t.Setenv(EnvFaviconEmoji, "🐧")
		assert.Contains(t, string(faviconHref()), "%F0%9F%90%A7")
	})
	t.Run("none falls back to /favicon.ico", func(t *testing.T) {
		t.Setenv(EnvFaviconEmoji, "none")
		assert.Empty(t, string(faviconHref()))
		body := renderPage(t, codeHandler, "/code")
		assert.Contains(t, body, `<link rel="icon" href="/favicon.ico">`)
		assert.NotContains(t, body, "image/svg+xml")
	})
	t.Run("markup in the value is escaped and long values are cut", func(t *testing.T) {
		t.Setenv(EnvFaviconEmoji, `</text><script>alert(1)</script>`)
		href := string(faviconHref())
		assert.NotContains(t, href, "<script")
		assert.NotContains(t, href, "%3Cscript") // cut to 8 characters, so no tag survives
		t.Setenv(EnvFaviconEmoji, "<b>")
		assert.Contains(t, string(faviconHref()), "%26lt%3Bb%26gt%3B") // &lt;b&gt;
	})
	t.Run("the page links the SVG icon with the ico as the fallback", func(t *testing.T) {
		t.Setenv(EnvFaviconEmoji, "")
		body := renderPage(t, codeHandler, "/code")
		assert.Contains(t, body, `rel="icon" type="image/svg+xml" href="data:image/svg&#43;xml,%3Csvg`) // html/template writes + as &#43;; browsers decode it
		assert.Contains(t, body, `rel="alternate icon" href="/favicon.ico"`)
	})
}

func TestCodePageLinksOnlyShowWhenSet(t *testing.T) {
	t.Run("GITLAB_URL=none hides its button (no dead link back to the same page)", func(t *testing.T) {
		t.Setenv(EnvGithubURL, "https://github.com/someone")
		t.Setenv(EnvGitlabURL, "none")
		body := renderPage(t, codeHandler, "/code")
		assert.Contains(t, body, `href="https://github.com/someone"`)
		assert.Contains(t, body, "GitHub")
		assert.NotContains(t, body, "GitLab")
		assert.NotContains(t, body, `href=""`)
	})
	t.Run("both show when both are set", func(t *testing.T) {
		t.Setenv(EnvGithubURL, "https://github.com/someone")
		t.Setenv(EnvGitlabURL, "https://gitlab.com/someone")
		body := renderPage(t, codeHandler, "/code")
		assert.Contains(t, body, `href="https://gitlab.com/someone"`)
		assert.Contains(t, body, "GitLab")
	})
}

func TestGitlabURL(t *testing.T) {
	t.Setenv(EnvGitlabURL, "")
	t.Setenv(EnvGitlabAccount, "")
	assert.Equal(t, "https://gitlab.com/lachie-gitlab", gitlabURL())

	t.Setenv(EnvGitlabAccount, "lachieggg")
	assert.Equal(t, "https://gitlab.com/lachieggg", gitlabURL())

	t.Setenv(EnvGitlabAccount, "bad/name")
	assert.Equal(t, "https://gitlab.com/lachie-gitlab", gitlabURL())

	t.Setenv(EnvGitlabURL, "https://example.com/me")
	assert.Equal(t, "https://example.com/me", gitlabURL())

	t.Setenv(EnvGitlabURL, "none")
	assert.Equal(t, "", gitlabURL())
}

func TestGithubURL(t *testing.T) {
	t.Setenv(EnvGithubURL, "")
	t.Setenv(EnvGithubAccount, "")
	assert.Equal(t, "https://github.com/lachieggg", githubURL())

	t.Setenv(EnvGithubAccount, "someone")
	assert.Equal(t, "https://github.com/someone", githubURL())

	t.Setenv(EnvGithubAccount, "bad/name")
	assert.Equal(t, "https://github.com/lachieggg", githubURL())

	t.Setenv(EnvGithubURL, "https://example.com/me")
	assert.Equal(t, "https://example.com/me", githubURL())

	t.Setenv(EnvGithubURL, "none")
	assert.Equal(t, "", githubURL())
}
