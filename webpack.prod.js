const path = require("path");
const merge = require("webpack-merge");
const common = require("./webpack.common");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const TerserPlugin = require("terser-webpack-plugin");
const OptimizeCssAssetsPlugin = require("optimize-css-assets-webpack-plugin");
const {CleanWebpackPlugin} = require("clean-webpack-plugin");

module.exports = merge(common, {
    mode: "production",
    output: {
        filename: "js/[name]-[contentHash].min.js",
        path: path.resolve(__dirname, "web/generated")
    },
    optimization: {
        minimizer: [
            new OptimizeCssAssetsPlugin(),
            new TerserPlugin(),
            new HtmlWebpackPlugin(
                {
                    chunks: ['app'],
                    template: "template/frontend/httpd/index.html",
                    filename: "index.html",
                    favicon: "./web/favicon.ico",
                    base: "/",
                    minify: {
                        removeAttributeQuotes: true,
                        collapseWhitespace: true,
                        removeComments: true
                    }
                }
            ),
        ]
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: "css/[name]-[contentHash].css"
        }),
        new CleanWebpackPlugin()
    ],
    module: {
        rules: [
            {
                test: /\.css$/,
                exclude: /node_modules/,
                use: [
                    MiniCssExtractPlugin.loader,
                    {
                        loader: "css-loader"
                    }
                ],
            },
        ],
    }
});
