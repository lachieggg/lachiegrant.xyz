const path = require('path');
const Dotenv = require('dotenv-webpack');

module.exports = {
    mode: 'production',
    entry: {
        'home': './src/js/home.js',
	'index':'./src/js/index.js'
    },
    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, 'public/js/generated'),
    },
    plugins: [
        new Dotenv()
    ]
};
