package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// cssPath and homeJSPath are computed once at startup by publishVersionedAsset,
// which content-hashes a static file and republishes it under a hash-named
// path, e.g. /css/app.a1b2c3d4e5f6.css. This is not just a nice-to-have:
// nginx sets `expires 1y` on .css/.js (nginx/nginx.conf.template), and
// Cloudflare honors that origin header per edge location worldwide - every
// one of Cloudflare's edge locations that has ever fetched a plain
// /css/app.css or /js/generated/home.js caches it for a full year
// independently. A query-string cache-buster (?v=hash) was ALSO observed
// serving stale content from at least one edge location despite being a
// "new" URL. A genuinely different PATH is the only thing that's a
// guaranteed-new cache key everywhere, immune to both of those failure
// modes, which is why this exists instead of either alternative.
var (
	cssPath        = "/css/app.css"          // overwritten at startup; safe fallback
	homeJSPath     = "/js/generated/home.js" // overwritten at startup; safe fallback
	logoPath       = "/images/rpi.png"       // overwritten at startup; safe fallback
	debianLogoPath = "/images/debian.png"    // overwritten at startup; safe fallback
)

// homeImageVersions maps a configured picture's filename (e.g. "beauty.jpg")
// to its published, content-hash-versioned URL - populated once at startup
// by publishHomeImages. See homeImage() in handlers.go for why this
// exists: a plain /images/<name> URL can get stuck serving a stale file
// from a Cloudflare edge cache for up to a year (expires 1y in nginx) even
// after the origin file changes underneath it - exactly what happened to
// the footer's rpi.png logo, which kept serving a ~309KB file for almost
// two days after it was replaced with a ~4KB one.
var homeImageVersions = map[string]string{}

// publishHomeImages content-hash-versions every image homeImage() might
// serve: the fallback placeholder plus every entry in PICTURES. Called
// once at startup (see main.go), not per-request - these are full-size
// photos, hashing them on every request would be wasteful.
func publishHomeImages() {
	names := append([]string{defaultHomeImage}, configuredPictures()...)
	for _, name := range names {
		if _, ok := homeImageVersions[name]; ok {
			continue // already published (duplicate entry in PICTURES)
		}
		v, err := publishVersionedAsset("public/images/" + name)
		if err != nil {
			logger.Printf("Warning: could not publish versioned image %s (falling back to /images/%s): %v", name, name, err)
			continue
		}
		homeImageVersions[name] = v
	}
}

// publishVersionedAsset hashes srcPath's content, writes a hash-named copy
// into the same directory, removes any older versioned copies of that same
// file left over from a previous run, and returns the public URL path to
// the new copy. srcPath must be under "public/" - that prefix is stripped
// to build the URL (public/css/app.css -> /css/..., public/js/generated/
// home.js -> /js/generated/...).
func publishVersionedAsset(srcPath string) (string, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", srcPath, err)
	}

	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])[:12]

	dir := filepath.Dir(srcPath)
	base := filepath.Base(srcPath)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)

	versionedName := fmt.Sprintf("%s.%s%s", stem, hash, ext)
	versionedFullPath := filepath.Join(dir, versionedName)

	if err := os.WriteFile(versionedFullPath, data, 0644); err != nil {
		return "", fmt.Errorf("writing %s: %w", versionedFullPath, err)
	}

	cleanupOldVersionedAssets(dir, stem, ext, versionedName)

	rel := strings.TrimPrefix(srcPath, "public/")
	urlDir := "/" + filepath.ToSlash(filepath.Dir(rel))
	return urlDir + "/" + versionedName, nil
}

// cleanupOldVersionedAssets removes previously-published hash-named copies
// of one specific file (from earlier process starts) so they don't
// accumulate on disk forever. Best-effort: logs and continues on any error
// rather than failing startup.
func cleanupOldVersionedAssets(dir, stem, ext, keep string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	prefix := stem + "."
	for _, e := range entries {
		name := e.Name()
		if name == keep || name == stem+ext {
			continue // keep the current version and the master file itself
		}
		if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, ext) {
			if err := os.Remove(filepath.Join(dir, name)); err != nil {
				logger.Printf("Warning: could not remove old versioned asset %s: %v", name, err)
			}
		}
	}
}
