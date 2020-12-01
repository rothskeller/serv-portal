import path from 'path'

module.exports = {
  root: 'client',
  proxy: {
    '/api': 'http://localhost:8100/',
    '/dl': 'http://localhost:8100/',
  },
}
