var path = require('path');

module.exports = {
    entry: './build/index.js',
    mode: 'development',
    output: {
        path: __dirname + '/build/',
        filename: 'shell.bundle.js',
        publicPath: './build/'
    },
    module: {
        rules: [
            { test: /\.css$/, use: ['style-loader', 'css-loader'] },
            { test: /\.png$/, use: 'file-loader' }
        ]
    }
};
