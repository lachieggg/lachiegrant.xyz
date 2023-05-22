const URL = 'images/'
const IMAGES = process.env.PICTURES.split(',');
const HOME_IMG = process.env.HOME_IMG;
const BLOG_IMG = process.env.BLOG_IMG;

const HOME_ID = "home-img"
const BLOG_ID = "blog-img"
const HOME_BTN = "home-img-btn"

setPicture(HOME_IMG, HOME_ID);
setPicture(BLOG_IMG, BLOG_ID)

function setPicture(name, id) {
    // Get the image URL from the name
    var imageUrl = URL + name

    try {
        var image = document.getElementById(id);
        image.src = imageUrl;
    } catch (err) {
        // no home image on page
        // skipping
    }
}

function randomHomePicture() {
    // Get a random index within the array length
    var randomIndex = Math.floor(Math.random() * IMAGES.length);
    var imageName = IMAGES[randomIndex];

    setPicture(imageName, HOME_ID)
}

document.getElementById(HOME_BTN).addEventListener('click', randomHomePicture);
