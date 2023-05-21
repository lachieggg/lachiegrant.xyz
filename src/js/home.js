var IMAGES_URL = process.env.IMAGES_URL + 'images/'

const images = process.env.PICTURES.split(',');

// Get a random index within the array length
var randomIndex = Math.floor(Math.random() * images.length);

var defaultImage = IMAGES_URL + images[randomIndex];

setPicture();

function setPicture()
{
    try {
        document.getElementById("home-img").src = defaultImage;
    } catch (err) {
        // no home image on page
        // skipping
    }
}