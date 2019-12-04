const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CleanWebpackPlugin = require('clean-webpack-plugin');

module.exports = {
  entry: './src/js/main.js',
  plugins: [
    new CleanWebpackPlugin(['dist']),  // clean dist
    new HtmlWebpackPlugin({
        filename: 'index.html',
        template: 'index.html',
        title: '吃豆子',
        inject: true,
        minify: {
            removeComments: true,
            collapseWhitespace: true
        }
    })
  ],
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist')
  }
};