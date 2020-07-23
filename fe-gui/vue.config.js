module.exports = {
    publicPath: '/app/',
    lintOnSave: false,
    runtimeCompiler: true,
    configureWebpack: {
        //Necessary to run npm link https://webpack.js.org/configuration/resolve/#resolve-symlinks
        resolve: {
            symlinks: false
        },
        devServer: {
            //To fix "Invalid Host header": https://github.com/gitpod-io/gitpod/issues/26
            disableHostCheck: true
        }
    },
}
