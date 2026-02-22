const URL = 'images/';
const IMAGES = process.env.PICTURES ? process.env.PICTURES.split(',') : [];
const HOME_ID = "home-img";
const BLOG_ID = "blog-img";
const BOOK_ID = "book-img";
const HOME_BTN = "home-img-btn";

// Set initial random images
setPicture(HOME_ID, "");
setPicture(BLOG_ID, "");
setPicture(BOOK_ID, "");

function setPicture(id, name) {
    if (!name) {
        name = getRandomImage();
    }
    if (!name) return;

    const imageUrl = URL + name;
    try {
        const image = document.getElementById(id);
        if (image) {
            image.src = imageUrl;
        }
    } catch (err) {
        console.warn("Could not set image for id:", id);
    }
}

function getRandomImage() {
    if (IMAGES.length === 0) return "";
    const randomIndex = Math.floor(Math.random() * IMAGES.length);
    return IMAGES[randomIndex].trim();
}

function setRandomImage(id) {
    const imageName = getRandomImage();
    if (imageName) {
        setPicture(id, imageName);
    }
}

function createRandomImageSetter(id) {
    return function () {
        setRandomImage(id)
    }
}

const homeButton = document.getElementById(HOME_BTN);
if (homeButton) {
    homeButton.addEventListener('click', createRandomImageSetter(HOME_ID));
}
