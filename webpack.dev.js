const path = require("path");
const {merge} = require("webpack-merge");
const common = require("./webpack.common");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const {CleanWebpackPlugin} = require("clean-webpack-plugin");

module.exports = merge(common, {
    mode: "development",
    output: {
        filename: "js/[name].js",
        path: path.resolve(__dirname, "web/generated")
    },
    plugins: [
        new HtmlWebpackPlugin(
            {
                chunks: ['app'],
                template: "template/frontend/httpd/index.html",
                filename: "index.html",
                favicon: "./web/favicon.ico",
                base: "/",
            }
        ),
        new CleanWebpackPlugin()
    ],
    module: {
        rules: [
            {
                test: /\.css$/,
                use: [
                    'vue-style-loader',
                    'css-loader'
                ]
            }
        ],
    },
    resolve: {
        alias: {
            vue: 'vue/dist/vue.js'
        }
    }
});
