/** Jest identity proxy for CSS Module classes; no browser styling is executed in tests. */
module.exports = new Proxy({}, {
  get(_target, key) {
    if (typeof key === 'string' && key !== '__esModule') {
      return key;
    }

    return undefined;
  },
});
