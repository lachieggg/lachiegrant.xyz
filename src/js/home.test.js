/**
 * @jest-environment jsdom
 */

// Mock process.env
process.env.PICTURES = 'localhost-1.png,localhost-2.png';

const { getRandomImage, setPicture, URL, IMAGES } = require('./home');

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
