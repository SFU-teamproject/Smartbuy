const { '@': alias } = require('./paths');

module.exports = function override(config) {
  config.resolve.alias = { ...config.resolve.alias, ...alias };
  return config;
};