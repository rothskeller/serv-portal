module.exports = {
    root: 'client',
    proxy: {
        '/api': 'http://localhost:8100/',
        '/dl': 'http://localhost:8100/',
    },
    esbuildTarget: ['es2020', 'safari11']
};