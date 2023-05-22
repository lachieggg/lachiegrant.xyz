const URL = 'images/'
const IMAGES = process.env.PICTURES.split(',');

setPicture();

function setPicture()
{
    // Get a random index within the array length
    var randomIndex = Math.floor(Math.random() * IMAGES.length);
    var imageUrl = URL + IMAGES[randomIndex];

    try {
        var image = document.getElementById("home-img");
        image.src = imageUrl;
    } catch (err) {
        // no home image on page
        // skipping
    }
}

export default setPicture;