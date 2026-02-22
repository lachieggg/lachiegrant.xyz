const URL = 'images/',
    HOME_ID = 'home-img',
    BLOG_ID = 'blog-img',
    BOOK_ID = 'book-img',
    HOME_BTN = 'home-img-btn',
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

if (IMAGES.length === 0) {
    console.error('PICTURES environment is empty or missing');
} else {
    // Set initial random images
    setRandomImage(HOME_ID);
    setRandomImage(BLOG_ID);
    setRandomImage(BOOK_ID);

    const homeButton = document.getElementById(HOME_BTN);
    if (homeButton) {
        homeButton.addEventListener('click', () => setRandomImage(HOME_ID));
    }
}

module.exports = {
    getRandomImage,
    setPicture,
    setRandomImage,
    URL,
    IMAGES
};
