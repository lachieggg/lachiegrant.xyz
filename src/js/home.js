const URL = 'images/';
const IMAGES = process.env.PICTURES.split(',');
const BLOG_IMG = process.env.BLOG_IMG;
const BOOK_IMG = process.env.BOOK_IMG;
const HOME_IMG = process.env.HOME_IMG;
const HOME_ID = "home-img";
const BLOG_ID = "blog-img";
const BOOK_ID = "book-img";
const HOME_BTN = "home-img-btn";

setPicture(HOME_ID,HOME_IMG);
setPicture(BLOG_ID,BLOG_IMG);
setPicture(BOOK_ID, BOOK_IMG);

// setPicture
// id   -> the id of the element to set 
// name -> the filename of the image
function setPicture(id, name) {
    if (name == "") {
        name = getRandomImage()
    }
    // Get the image URL from the name
    var imageUrl = URL + name;

    try {
        var image = document.getElementById(id);
        image.src = imageUrl;
    } catch (err) {
        // no image on page
        // skipping
    }
}

// getRandomImage returns a random image
// filename.
function getRandomImage() {
    var randomIndex = Math.floor(Math.random() * IMAGES.length);
    return IMAGES[randomIndex]
}

function setRandomImage(id) {
    imageName = getRandomImage()

    // Get a random index within the array length
    setPicture(HOME_ID,);
    if(imageName == "") {
        setRandomImage(id);
    }
}

// createRandomImageSetter creates an image
// setter function for an element id 
function createRandomImageSetter(id) {
    return function() {
        setRandomImage(id) 
    }
}

var homeButton = document.getElementById(HOME_BTN);
if(homeButton != null) {
    homeButton.addEventListener('click', createRandomImageSetter(HOME_ID));
}
