module.exports = {
  chainWebpack: (config) => {
    config.plugins.delete('prefetch')
  },
  configureWebpack: {
    resolve: {
      alias: {
        '@': __dirname + '/client',
      }
    },
    entry: {
      app: './client/main.js',
    },
  },
  devServer: {
    proxy: 'http://localhost:8100'
  },
}
