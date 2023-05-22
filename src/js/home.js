const URL = 'images/'
const IMAGES = process.env.PICTURES.split(',');
const DEFAULT = process.env.DEFAULT;

setPicture(DEFAULT);

function setPicture(name) {
    // Get the image URL from the name
    var imageUrl = URL + name

    try {
        var image = document.getElementById("home-img");
        image.src = imageUrl;
    } catch (err) {
        // no home image on page
        // skipping
    }
}

function randomPicture() {
    // Get a random index within the array length
    var randomIndex = Math.floor(Math.random() * IMAGES.length);
    var imageName = IMAGES[randomIndex];

    setPicture(imageName)
}

document.getElementById('home-img-btn').addEventListener('click', randomPicture);
