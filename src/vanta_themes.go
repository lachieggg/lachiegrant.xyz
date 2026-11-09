package main

import (
	"encoding/json"
	"html/template"
	"os"
)

// EnvVantaTheme selects which home-page background animation renders.
// Swapping themes is a one env var change + `docker compose restart app` -
// no code edit needed. Adding a new theme is one entry in vantaThemes.
const EnvVantaTheme = "VANTA_THEME"
const EnvHeroMode = "HERO_MODE"

const defaultVantaThemeName = "cells"
const defaultHeroMode = "animated"

// activeHeroMode reads HERO_MODE (default "animated"): "animated" renders
// the Vanta.js background (see activeVantaTheme for which effect); "static"
// renders the original plain version - just the image, no WebGL, nothing
// loaded from a CDN. This is the actual on/off switch for the animation,
// independent of which effect is selected - useful for browsers that
// restrict WebGL for fingerprinting protection (e.g. Brave Shields, which
// farbles WebGL extensions by design - not a bug to work around, just a
// privacy choice some visitors make) or for anyone who prefers the plain
// look.
func activeHeroMode() string {
	mode := os.Getenv(EnvHeroMode)
	if mode != "static" {
		return defaultHeroMode
	}
	return "static"
}

// VantaTheme describes one swappable Vanta.js background effect: which
// rendering engine it needs (Vanta effects are built on either three.js or
// p5.js, not always the same one), which effect script to load, the global
// init function it exposes, and its effect-specific parameters.
type VantaTheme struct {
	Label       string
	EngineJSURL string
	EffectJSURL string
	InitFunc    string
	Params      map[string]interface{}
}

// vantaThemes is the whole registry - every theme available to swap to.
// All sci-fi/nature-tech leaning, MIT licensed (github.com/tengbao/vanta).
var vantaThemes = map[string]VantaTheme{
	"cells": {
		Label:       "Cells - organic Voronoi cell pattern",
		EngineJSURL: "https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js",
		EffectJSURL: "https://cdn.jsdelivr.net/npm/vanta@0.5.24/dist/vanta.cells.min.js",
		InitFunc:    "CELLS",
		Params: map[string]interface{}{
			"mouseControls": true, "touchControls": true, "gyroControls": false,
			"minHeight": 200.0, "minWidth": 200.0, "scale": 1.0, "scaleMobile": 1.0,
			"color1": 0x0d3b2e, "color2": 0x1fae82, "size": 1.4, "speed": 0.9,
		},
	},
	"topology": {
		Label:       "Topology - flowing wireframe terrain contours",
		EngineJSURL: "https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js",
		EffectJSURL: "https://cdn.jsdelivr.net/npm/vanta@0.5.24/dist/vanta.topology.min.js",
		InitFunc:    "TOPOLOGY",
		Params: map[string]interface{}{
			"mouseControls": true, "touchControls": true, "gyroControls": false,
			"minHeight": 200.0, "minWidth": 200.0, "scale": 1.0, "scaleMobile": 1.0,
			"color": 0x1fae82, "backgroundColor": 0x0a1f1a,
		},
	},
	"trunk": {
		Label:       "Trunk - concentric growth-ring / radar-sweep rings",
		EngineJSURL: "https://cdnjs.cloudflare.com/ajax/libs/p5.js/1.9.4/p5.min.js",
		EffectJSURL: "https://cdn.jsdelivr.net/npm/vanta@0.5.24/dist/vanta.trunk.min.js",
		InitFunc:    "TRUNK",
		Params: map[string]interface{}{
			"mouseControls": true, "touchControls": true, "gyroControls": false,
			"minHeight": 200.0, "minWidth": 200.0, "scale": 1.0, "scaleMobile": 1.0,
			"backgroundColor": 0x0a1a17, "color": 0x1fae82, "spacing": 8.0, "chaos": 3.0,
		},
	},
}

// activeVantaTheme reads VANTA_THEME (default "cells"), falling back to the
// default for an unrecognized name rather than failing the page.
func activeVantaTheme() (string, VantaTheme) {
	name := os.Getenv(EnvVantaTheme)
	if name == "" {
		name = defaultVantaThemeName
	}
	if t, ok := vantaThemes[name]; ok {
		return name, t
	}
	return defaultVantaThemeName, vantaThemes[defaultVantaThemeName]
}

// vantaParamsJSON renders the theme's params (plus the fixed "el" target)
// as a JSON object literal, typed template.JS so html/template embeds it
// verbatim inside a <script> body. This only ever gets called inside a
// <script> context, never inside an HTML attribute - the nested-quote
// escaping bug this project hit earlier was specifically about attribute
// contexts, which this avoids entirely.
func vantaParamsJSON(t VantaTheme) template.JS {
	params := make(map[string]interface{}, len(t.Params)+1)
	for k, v := range t.Params {
		params[k] = v
	}
	params["el"] = "#vanta-bg"

	data, err := json.Marshal(params)
	if err != nil {
		return template.JS("{}")
	}
	return template.JS(data)
}
