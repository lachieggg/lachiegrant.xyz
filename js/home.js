var IMAGES_URL = process.env.IMAGES_URL + 'images/'

const images = process.env.PICTURES.split(',');

var defaultImage = IMAGES_URL + images[0]

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