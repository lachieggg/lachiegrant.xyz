package main

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	FeatureEnabled  = "1"
	FeatureDisabled = "0"

	EnvEnableBlog      = "ENABLE_BLOG"
	EnvEnableBookmarks = "ENABLE_BOOKMARKS"
	EnvGithubURL       = "GITHUB_URL"     // full address; overrides GITHUB_ACCOUNT, "none" hides the button
	EnvGithubAccount   = "GITHUB_ACCOUNT" // github.com username
	EnvGitlabURL       = "GITLAB_URL"     // full address; overrides GITLAB_ACCOUNT, "none" hides the button
	EnvGitlabAccount   = "GITLAB_ACCOUNT" // gitlab.com username, e.g. lachieggg
	EnvResumePath      = "RESUME_PATH"

	// The GitLab account shown when neither GITLAB_URL nor GITLAB_ACCOUNT is
	// set. The other account is lachieggg: set GITLAB_ACCOUNT=lachieggg to switch.
	defaultGitlabAccount = "lachie-gitlab"
	defaultGithubAccount = "lachieggg" // the only GitHub account

	TmplHome      = "home.html"
	TmplBookmarks = "bookmarks.html"
	TmplBlogIndex = "blog.html"
	TmplCode      = "code.html"
)

// parseTemplates loads and parses all HTML templates from the templates directory.
// Returns nil and writes an error response if template parsing fails.
// NOTE: Templates are currently re-parsed on each request. For production,
// consider parsing once at startup and caching the result.
func parseTemplates(w http.ResponseWriter) *template.Template {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, fmt.Sprintf("Operating system error: %v", err), http.StatusInternalServerError)
		return nil
	}

	// Parse standard page templates
	standardFolder := filepath.Join(wd, "src", "templates", "*.html")
	t, err := template.ParseGlob(standardFolder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing templates: %v", err), http.StatusInternalServerError)
		return nil
	}

	// Parse blog post templates
	blogFolder := filepath.Join(wd, "src", "templates", "blog", "*.html")
	t, err = t.ParseGlob(blogFolder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing templates: %v", err), http.StatusInternalServerError)
		return nil
	}

	return t
}

// PageData holds data passed to HTML templates.
type PageData struct {
	EnableBlog      bool
	EnableBookmarks bool
	CSSPath         string
	HomeJSPath      string
	LogoPath        string

	// Footer "hosted from" line and tab icon (HOSTED_FROM, FAVICON_EMOJI - see branding.go).
	HostedLabel   string
	HostedLogo    string
	HostedLogoAlt string
	FaviconHref   template.URL

	// Vanta background: site-wide, not home-page-specific (see header.html's
	// #vanta-bg mount and scripts.html's init script, both gated on
	// HeroMode) - every page gets the same animated/static choice and theme
	// so the theme-cycle toggle (src/js/home.js) behaves consistently no
	// matter which page it's clicked from, rather than only working on "/".
	HeroMode        string // "animated" or "static" - see HERO_MODE
	VantaName       string
	VantaEngineURL  string
	VantaEffectURL  string
	VantaInitFunc   template.JS // bare identifier (e.g. CELLS) - must not be quoted
	VantaParamsJSON template.JS

	Data interface{}
}

// getPageData returns a PageData object with common feature flags and the
// active Vanta background theme (see src/vanta_themes.go) populated.
func getPageData(data interface{}) PageData {
	blog := os.Getenv(EnvEnableBlog) == FeatureEnabled
	bookmarks := os.Getenv(EnvEnableBookmarks) == FeatureEnabled
	name, theme := activeVantaTheme()
	hosting := hostedFrom()
	return PageData{
		EnableBlog:      blog,
		EnableBookmarks: bookmarks,
		CSSPath:         cssPath,
		HomeJSPath:      homeJSPath,
		LogoPath:        logoPath,
		HostedLabel:     hosting.Label,
		HostedLogo:      hosting.Logo,
		HostedLogoAlt:   hosting.LogoAlt,
		FaviconHref:     faviconHref(),
		HeroMode:        activeHeroMode(),
		VantaName:       name,
		VantaEngineURL:  theme.EngineJSURL,
		VantaEffectURL:  theme.EffectJSURL,
		VantaInitFunc:   template.JS(theme.InitFunc),
		VantaParamsJSON: vantaParamsJSON(theme),
		Data:            data,
	}
}

// HomeData is passed as PageData.Data for the home page: just the hero
// image now that the Vanta background itself is site-wide (see PageData).
type HomeData struct {
	HomeImage string // URL, e.g. "/images/beauty.a1b2c3.jpg" - see homeImage()
}

const EnvPictures = "PICTURES"
const defaultHomeImage = "sample.png"

// configuredPictures returns the parsed, trimmed PICTURES list (comma
// separated filenames under public/images/), or nil if unset/empty.
func configuredPictures() []string {
	raw := os.Getenv(EnvPictures)
	var pics []string
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			pics = append(pics, p)
		}
	}
	return pics
}

// pickHomeImageName picks which configured picture to use for this
// request - a random pick among PICTURES entries, or defaultHomeImage if
// none are configured (e.g. a fresh clone with no .env yet).
func pickHomeImageName() string {
	pics := configuredPictures()
	if len(pics) == 0 {
		return defaultHomeImage
	}
	return pics[rand.Intn(len(pics))]
}

// homeImage picks the hero image server-side so its real path is already in
// the initial HTML response. This used to be a client-side-only choice
// (home.js reading process.env.PICTURES, baked in at webpack build time) -
// the HTML always shipped a hardcoded placeholder src (images/sample.png,
// a ~173-byte near-blank PNG) that JS then swapped out once the deferred JS
// bundle had fetched, parsed and run. On a warm browser cache that swap is
// instant and invisible; on a cold cache/hard refresh (no local cache to
// serve the real photo from immediately) that gap is visible as a blank
// placeholder box flashing before the real photo loads. Picking it here
// means the real photo starts downloading the moment the HTML is parsed,
// in parallel with CSS/JS, instead of waiting on a JS round-trip.
//
// The returned path is the content-hash-versioned URL published by
// publishHomeImages() at startup (see assets.go), not a plain
// /images/<name> path - a plain path is exactly what let the footer logo
// (rpi.png) get stuck serving a stale file from a Cloudflare edge cache for
// nearly two days after the origin file changed (expires 1y in nginx, and a
// same-path swap doesn't invalidate an edge's existing cache entry). Falls
// back to a plain path only if that image somehow wasn't published -
// should not normally happen, but must never break the homepage over it.
func homeImage() string {
	name := pickHomeImageName()
	if v, ok := homeImageVersions[name]; ok {
		return v
	}
	return "/images/" + name
}

func getHomeData() HomeData {
	return HomeData{HomeImage: homeImage()}
}

// indexHandler serves the root path. It renders the home page for "/" and
// returns 404 for any other path (since "/" matches all unhandled routes).
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		logger.Printf("404: %s not found (matched index catch-all)", r.URL.Path)
		renderNotFound(w, r)
		return
	}
	homeHandler(w, r)
}

// homeHandler renders the home page template.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	t := parseTemplates(w)
	if t == nil {
		return
	}

	if err := t.ExecuteTemplate(w, TmplHome, getPageData(getHomeData())); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

// CodeData is passed as PageData.Data for the code page: where each
// platform button links to. GitHub's logo usage is explicitly permitted
// for exactly this - "use...the Invertocat logo as a social button to link
// to your GitHub profile" (github.com/logos). GitLab's own trademark
// guidelines are less explicit about this exact case, but as a GitLab
// employee linking to a personal profile the risk is minimal - see
// conversation/commit context.
type CodeData struct {
	GithubURL string
	GitlabURL string
}

// profileURL resolves a code-host button's address, identically for GitHub
// and GitLab: the full-address setting if set ("none" hides the button), else
// <host>/<account setting>, else <host>/<default account>.
func profileURL(urlEnv, accountEnv, host, defaultAccount string) string {
	if u := strings.TrimSpace(os.Getenv(urlEnv)); u != "" {
		if strings.EqualFold(u, "none") {
			return ""
		}
		return u
	}
	account := strings.Trim(strings.TrimSpace(os.Getenv(accountEnv)), "/")
	if account == "" || strings.ContainsAny(account, " /?#") {
		account = defaultAccount
	}
	return "https://" + host + "/" + account
}

func githubURL() string {
	return profileURL(EnvGithubURL, EnvGithubAccount, "github.com", defaultGithubAccount)
}

func gitlabURL() string {
	return profileURL(EnvGitlabURL, EnvGitlabAccount, "gitlab.com", defaultGitlabAccount)
}

func getCodeData() CodeData {
	return CodeData{
		GithubURL: githubURL(),
		GitlabURL: gitlabURL(),
	}
}

// codeHandler renders the code page: a prototype for choosing between two
// platforms (GitHub, GitLab - the latter added since the site owner is a
// GitLab employee too) rather than redirecting straight to one.
func codeHandler(w http.ResponseWriter, r *http.Request) {
	t := parseTemplates(w)
	if t == nil {
		return
	}

	if err := t.ExecuteTemplate(w, TmplCode, getPageData(getCodeData())); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

// bookmarksHandler renders the bookmarks page template.
func bookmarksHandler(w http.ResponseWriter, r *http.Request) {
	if os.Getenv(EnvEnableBookmarks) != FeatureEnabled {
		logger.Printf("404: Bookmarks disabled via feature flag")
		renderNotFound(w, r)
		return
	}

	t := parseTemplates(w)
	if t == nil {
		return
	}

	if err := t.ExecuteTemplate(w, TmplBookmarks, getPageData(nil)); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

// resumeHandler serves the resume PDF file.
func resumeHandler(w http.ResponseWriter, r *http.Request) {
	resumePath := os.Getenv(EnvResumePath)
	if resumePath == "" {
		logger.Printf("Error: RESUME_PATH environment variable not set")
		http.Error(w, "Resume path not configured", http.StatusInternalServerError)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		logger.Printf("Error getting working directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Clean and resolve the absolute path relative to the 'public' directory
	publicDir := filepath.Join(wd, "public")
	absPath := filepath.Join(publicDir, filepath.Clean(resumePath))

	// Security Check: Ensure the resolved path is still within the 'public' directory
	if !strings.HasPrefix(absPath, publicDir) {
		logger.Printf("Security Warning: Unauthorized access attempt to %s", absPath)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		logger.Printf("404: Resume file not found at %s", absPath)
		renderNotFound(w, r)
		return
	}

	logger.Printf("Serving resume from: %s", absPath)
	http.ServeFile(w, r, absPath)
}

// renderNotFound renders the site-themed 404 page (header, nav, footer, theme
// script) with a 404 status. It renders to a buffer first so a template
// problem falls back to the plain http.NotFound instead of a half-written page.
func renderNotFound(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	t, err := template.ParseGlob(filepath.Join(wd, "src", "templates", "*.html"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "404.html", getPageData(nil)); err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(buf.Bytes())
}
