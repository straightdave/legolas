const path = require('path')

module.exports = {
    entry: './app/src/index.js',

    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'public/dist')
    },

    module: {
        rules: [
            {
                test: /\.vue$/,
                use: ['vue-loader']
            }
        ]
    }
}
