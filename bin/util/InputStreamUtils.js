const JSONStream = require("JSONStream");

function readJsonInputStream(jsonPath) {
  return process.openStdin().pipe(JSONStream.parse(jsonPath));
}

module.exports = { readJsonInputStream };
