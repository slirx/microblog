const path = require("path");
const {VueLoaderPlugin} = require("vue-loader");

module.exports = {
    entry: {
        app: "./web/js/app.js"
    },
    output: {
        filename: "js/[name]-[contentHash].min.js",
        path: path.resolve(__dirname, "web/generated")
    },
    plugins: [
        new VueLoaderPlugin()
    ],
    module: {
        rules: [
            {
                test: /\.js$/,
                exclude: /node_modules/,
                loader: "babel-loader"
            },
            {
                test: /\.vue$/,
                exclude: /node_modules/,
                loader: "vue-loader"
            },
            {
                test: /\.html$/,
                use: [
                    {
                        loader: "html-loader"
                    }
                ]
            },
            {
                test: /\.(png|jpg|gif)$/i,
                use: [
                    {
                        loader: "file-loader",
                        options: {
                            name: this.mode === 'production' ? "[name]-[contenthash].[ext]" : "[name].[ext]",
                            outputPath: "/images",
                            esModule: false
                        }
                    }
                ],
            },
            {
                test: /\.(ico)$/i,
                use: [
                    {
                        loader: "file-loader",
                        options: {
                            name: "[name].[ext]",
                            outputPath: "/",
                            esModule: false
                        }
                    }
                ],
            },
        ],
    }
};
