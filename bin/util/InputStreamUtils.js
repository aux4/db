const JSONStream = require("JSONStream");
const {Transform} = require("stream");

function readJsonInputStream(jsonPath) {
  return process.openStdin().pipe(JSONStream.parse(jsonPath)).pipe(new SplitterTransform());
}

class SplitterTransform extends Transform {
  constructor() {
    super({
      writableObjectMode: true,
      readableObjectMode: true
    });
  }

  _transform(value, encoding, callback) {
    if (Array.isArray(value)) {
      value.forEach(item => this.push(item));
      callback();
      return;
    }
    callback(null, value);
  }
}

module.exports = { readJsonInputStream };
