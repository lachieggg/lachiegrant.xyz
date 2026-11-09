/**
 * @jest-environment jsdom
 */

// Mock process.env
process.env.PICTURES = 'localhost-1.png,localhost-2.png';

const {
    getRandomImage,
    setPicture,
    URL,
    IMAGES,
    THEME_KEY,
    vantaAvailable,
    availablePresets,
    currentPreset,
    applyPreset,
    cyclePreset,
    initThemeToggle,
} = require('./home');

describe('home.js front-end generic testing logic', () => {
    beforeEach(() => {
        // Clear DOM and reset it to a generic testing state
        document.body.innerHTML = `
            <img id="test-img" src="" alt="test img">
        `;
    });

    test('getRandomImage securely retrieves a mapped image', () => {
        // If there are images provided logically, we expect it to return one matching a string.
        const returnedImg = getRandomImage();

        if (IMAGES.length > 0) {
            expect(typeof returnedImg).toBe('string');
            expect(IMAGES.map(i => i.trim())).toContain(returnedImg);
        } else {
            expect(returnedImg).toBe('');
        }
    });

    test('setPicture updates generic DOM identifier string cleanly', () => {
        const testId = "test-img";
        const dummyImageString = "localhost-mock.png";

        // Assert securely set
        setPicture(testId, dummyImageString);

        // Fetch element back from DOM logically
        const imgElement = document.getElementById(testId);

        // JSDOM will natively prepend localhost/ or URL matching. 
        // We ensure whatever the domain resolves to naturally, it natively appended the correct generic URL + mocked image naming convention securely
        expect(imgElement.src).toContain(`${URL}${dummyImageString}`);
    });

    test('setPicture gracefully skips if provided mismatched ID', () => {
        // Should not throw or crash on missing element
        expect(() => {
            setPicture("completely-missing-element", "mock.png");
        }).not.toThrow();
    });
});

describe('theme toggle (HERO_MODE=static - no vanta on the page)', () => {
    beforeEach(() => {
        localStorage.clear();
        document.documentElement.removeAttribute('data-theme');
        document.documentElement.removeAttribute('data-theme-preset');
        document.body.innerHTML = '<button id="theme-toggle"></button>';
        delete window.__vantaInit;
        delete window.__vantaDestroy;
        delete window.__vantaEffect;
    });

    test('vantaAvailable is false and the cycle degrades to 2-way light/dark', () => {
        expect(vantaAvailable()).toBe(false);
        expect(availablePresets()).toEqual(['light', 'dark']);
    });

    test('currentPreset defaults to light when no attribute is set', () => {
        expect(currentPreset()).toBe('light');
    });

    test('applyPreset sets and clears the data-theme attribute', () => {
        applyPreset('dark');
        expect(document.documentElement.getAttribute('data-theme')).toBe('dark');

        applyPreset('light');
        expect(document.documentElement.hasAttribute('data-theme')).toBe(false);
    });

    test('cyclePreset flips between light and dark and persists it to localStorage', () => {
        expect(cyclePreset()).toBe('dark');
        expect(localStorage.getItem(THEME_KEY)).toBe('dark');

        expect(cyclePreset()).toBe('light');
        expect(localStorage.getItem(THEME_KEY)).toBe('light');
    });

    test('clicking the toggle button cycles the theme (simulates the real DOM wiring)', () => {
        const btn = document.getElementById('theme-toggle');
        initThemeToggle();

        expect(currentPreset()).toBe('light');
        btn.click();
        expect(currentPreset()).toBe('dark');
        expect(btn.getAttribute('aria-pressed')).toBe('true');

        btn.click();
        expect(currentPreset()).toBe('light');
        expect(btn.getAttribute('aria-pressed')).toBe('false');
    });
});

describe('theme toggle (HERO_MODE=animated - vanta present on the page)', () => {
    beforeEach(() => {
        localStorage.clear();
        document.documentElement.removeAttribute('data-theme');
        document.documentElement.removeAttribute('data-theme-preset');
        document.body.innerHTML = '<button id="theme-toggle"></button>';
        window.__vantaEffect = null;
        window.__vantaInit = jest.fn(() => {
            window.__vantaEffect = {};
            return window.__vantaEffect;
        });
        window.__vantaDestroy = jest.fn(() => {
            window.__vantaEffect = null;
        });
    });

    afterEach(() => {
        delete window.__vantaInit;
        delete window.__vantaDestroy;
        delete window.__vantaEffect;
    });

    test('vantaAvailable is true and the cycle includes cells', () => {
        expect(vantaAvailable()).toBe(true);
        expect(availablePresets()).toEqual(['light', 'dark', 'cells']);
    });

    test('cyclePreset goes cells -> light -> dark -> cells, starting/stopping vanta', () => {
        // With nothing stored yet, the page has already auto-started vanta
        // by the time this file's own init runs (see home.html) - so the
        // live starting point is 'cells', not 'light'.
        expect(currentPreset()).toBe('cells');

        expect(cyclePreset()).toBe('light');
        expect(document.documentElement.hasAttribute('data-theme')).toBe(false);
        expect(window.__vantaDestroy).toHaveBeenCalledTimes(1);

        expect(cyclePreset()).toBe('dark');
        expect(document.documentElement.getAttribute('data-theme')).toBe('dark');
        expect(window.__vantaInit).not.toHaveBeenCalled();

        expect(cyclePreset()).toBe('cells');
        expect(document.documentElement.getAttribute('data-theme')).toBe('dark');
        expect(document.documentElement.getAttribute('data-theme-preset')).toBe('cells');
        expect(window.__vantaInit).toHaveBeenCalledTimes(1);
    });

    test('clicking the toggle button cycles through all three presets', () => {
        const btn = document.getElementById('theme-toggle');
        initThemeToggle();

        expect(currentPreset()).toBe('cells');
        btn.click();
        expect(currentPreset()).toBe('light');
        btn.click();
        expect(currentPreset()).toBe('dark');
        btn.click();
        expect(currentPreset()).toBe('cells');
    });
});

describe('theme toggle (WebGL unavailable, e.g. Brave farbling)', () => {
    beforeEach(() => {
        localStorage.clear();
        document.documentElement.removeAttribute('data-theme');
        document.documentElement.removeAttribute('data-theme-preset');
        document.body.innerHTML = '<button id="theme-toggle"></button>';
        delete window.__vantaBroken;
        window.__vantaDestroy = jest.fn();
    });

    afterEach(() => {
        delete window.__vantaInit;
        delete window.__vantaDestroy;
        delete window.__vantaBroken;
        delete window.__vantaEffect;
    });

    test('applyPreset("cells") falls back to light when __vantaInit returns null, and reports that as the applied preset', () => {
        window.__vantaInit = jest.fn(() => null); // home.html's own try/catch already failed and returned null
        expect(applyPreset('cells')).toBe('light');
        expect(document.documentElement.hasAttribute('data-theme')).toBe(false);
        expect(document.documentElement.getAttribute('data-theme-preset')).toBe('light');
        expect(window.__vantaBroken).toBe(true);
    });

    test('applyPreset("cells") falls back to light when __vantaInit throws synchronously', () => {
        window.__vantaInit = jest.fn(() => { throw new Error('Error creating WebGL context.'); });
        expect(applyPreset('cells')).toBe('light');
        expect(document.documentElement.hasAttribute('data-theme')).toBe(false);
        expect(window.__vantaBroken).toBe(true);
    });

    test('once broken, cells drops out of the cycle entirely - no repeated failed attempts', () => {
        window.__vantaInit = jest.fn(() => null);
        localStorage.setItem(THEME_KEY, 'dark'); // one step before cells in the cycle

        expect(cyclePreset()).toBe('light'); // tried 'cells', it failed, fell back
        expect(window.__vantaInit).toHaveBeenCalledTimes(1);

        expect(availablePresets()).toEqual(['light', 'dark']);
        expect(cyclePreset()).toBe('dark');
        expect(cyclePreset()).toBe('light');
        // Never retried __vantaInit once broken.
        expect(window.__vantaInit).toHaveBeenCalledTimes(1);
    });

    test('clicking the toggle degrades to a plain binary light/dark cycle after a failed attempt', () => {
        window.__vantaInit = jest.fn(() => null);
        localStorage.setItem(THEME_KEY, 'dark'); // one click away from actually attempting cells
        const btn = document.getElementById('theme-toggle');
        initThemeToggle();

        expect(currentPreset()).toBe('dark');
        btn.click(); // tries to advance to 'cells', fails, falls back to 'light'
        expect(currentPreset()).toBe('light');
        expect(btn.getAttribute('aria-label')).toBe('Switch to dark theme');
        btn.click();
        expect(currentPreset()).toBe('dark');
        btn.click();
        expect(currentPreset()).toBe('light');
    });
});
