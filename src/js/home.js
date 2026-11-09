const URL = 'images/',
    BLOG_ID = 'blog-img',
    BOOK_ID = 'book-img',
    IMAGES = process.env.PICTURES ? process.env.PICTURES.split(',') : [];

function getRandomImage() {
    if (!IMAGES.length) return '';
    return IMAGES[Math.floor(Math.random() * IMAGES.length)].trim();
}

function setPicture(id, name) {
    try {
        const image = document.getElementById(id);
        if (image) image.src = URL + name;
    } catch (err) {
        console.warn('Could not set image for id:', id);
    }
}

function setRandomImage(id) {
    setPicture(id, getRandomImage());
}

// The home page's hero image (#home-img) is now picked server-side (see
// homeImage() in src/handlers.go) and rendered directly into the initial
// HTML response, so it starts downloading immediately instead of behind a
// placeholder (images/sample.png) that got swapped out once this bundle
// loaded - that gap was visible as a blank placeholder flash on a cold
// cache/hard refresh. Re-randomizing it here would either bring that flash
// back or, on a warm cache, flicker to a different photo right after the
// page has already painted. blog-img/book-img (bookmarks.html) don't have
// that server-side treatment yet, so they're still picked client-side.
if (IMAGES.length === 0) {
    console.error('PICTURES environment is empty or missing');
} else {
    setRandomImage(BLOG_ID);
    setRandomImage(BOOK_ID);
}

const THEME_KEY = 'theme';

// vantaAvailable checks for window.__vantaInit, which home.html only
// defines when HERO_MODE=animated (see src/vanta_themes.go). With
// HERO_MODE=static there's no vanta canvas to cycle through at all, so the
// cycle degrades gracefully to a plain 2-way light/dark toggle. It also
// respects window.__vantaBroken, set (by home.html or the catch block
// below) the first time WebGL actually fails to start - e.g. Brave's
// farbling fingerprinting protection refuses to create a WebGL context at
// all. Once that happens for this page load, 'cells' drops out of the
// cycle entirely rather than being retried and failing on every click.
function vantaAvailable() {
    return typeof window !== 'undefined' && typeof window.__vantaInit === 'function' && !window.__vantaBroken;
}

// availablePresets is the ordered cycle. More presets (other Vanta themes,
// etc.) can be added here in future without touching the rest of this file.
function availablePresets() {
    return vantaAvailable() ? ['light', 'dark', 'cells'] : ['light', 'dark'];
}

// currentPreset prefers the stored value: 'dark' and 'cells' both render
// data-theme="dark", so they can't be told apart from the DOM alone once
// something has been explicitly chosen. With nothing stored, the page is
// showing whatever HERO_MODE/VANTA_THEME rendered by default - an animated
// hero defaults to 'cells' (see header.html - it only suppresses vanta when
// the stored value is explicitly 'light' or 'dark'), a static hero falls
// back to whatever data-theme system preference produced.
function currentPreset() {
    try {
        const stored = localStorage.getItem(THEME_KEY);
        if (stored && availablePresets().includes(stored)) return stored;
    } catch (err) {}
    if (vantaAvailable()) return 'cells';
    return document.documentElement.getAttribute('data-theme') === 'dark' ? 'dark' : 'light';
}

// applyPreset applies a preset and returns the preset that actually ended
// up applied - normally the same as requested, but 'cells' falls back to
// 'light' if vanta fails to start (see vantaAvailable above). Callers that
// persist the result must persist the RETURN VALUE, not the argument, so a
// failed 'cells' attempt is never stored as if it had worked.
function applyPreset(preset) {
    if (preset === 'cells') {
        let effect = null;
        if (vantaAvailable()) {
            try {
                effect = window.__vantaInit();
            } catch (err) {
                console.warn('Vanta failed to start:', err);
                effect = null;
            }
            if (!effect) window.__vantaBroken = true;
        }
        if (!effect) {
            return applyPreset('light');
        }
        document.documentElement.setAttribute('data-theme', 'dark');
        document.documentElement.setAttribute('data-theme-preset', 'cells');
        return 'cells';
    }

    if (preset === 'dark') {
        document.documentElement.setAttribute('data-theme', 'dark');
    } else {
        document.documentElement.removeAttribute('data-theme');
    }
    // Lets CSS distinguish the 'cells' state from plain 'dark', since both
    // set the same data-theme value above.
    document.documentElement.setAttribute('data-theme-preset', preset);

    if (typeof window !== 'undefined' && typeof window.__vantaDestroy === 'function') {
        window.__vantaDestroy();
    }
    return preset;
}

function storePreset(preset) {
    try {
        localStorage.setItem(THEME_KEY, preset);
    } catch (err) {
        console.warn('Could not save theme preference:', err);
    }
}

// cyclePreset advances to the next preset in availablePresets(), applies
// and persists it, and returns the new value - kept separate from DOM event
// wiring so it can be tested without simulating a click.
function cyclePreset() {
    const presets = availablePresets();
    const idx = presets.indexOf(currentPreset());
    const next = presets[(idx + 1) % presets.length];
    const applied = applyPreset(next);
    storePreset(applied);
    return applied;
}

// initThemeToggle wires the corner swatch button (present on every page -
// see header.html). The initial preset is already applied before this file
// runs - light/dark chrome by header.html's blocking inline script, 'cells'
// by home.html's vanta auto-start - so this only handles cycling further.
function initThemeToggle() {
    const btn = document.getElementById('theme-toggle');
    if (!btn) return;

    function render() {
        const preset = currentPreset();
        document.documentElement.setAttribute('data-theme-preset', preset);
        btn.setAttribute('aria-pressed', String(preset !== 'light'));
        const presets = availablePresets();
        const next = presets[(presets.indexOf(preset) + 1) % presets.length];
        btn.setAttribute('aria-label', 'Switch to ' + next + ' theme');
    }

    btn.addEventListener('click', () => {
        cyclePreset();
        render();
    });

    render();
}

if (typeof document !== 'undefined') {
    initThemeToggle();
}

module.exports = {
    getRandomImage,
    setPicture,
    setRandomImage,
    URL,
    IMAGES,
    THEME_KEY,
    vantaAvailable,
    availablePresets,
    currentPreset,
    applyPreset,
    storePreset,
    cyclePreset,
    initThemeToggle
};
