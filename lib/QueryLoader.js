const fs = require("fs");

async function getQuery(file, query) {
  if (!file && !query) {
    throw new Error("file or query must be specified");
  }

  if (file && query) {
    throw new Error("file and query cannot be specified at the same time");
  }

  if (query) {
    return query;
  }

  return fs.readFileSync(file, { encoding: "utf8" });
}

module.exports = { getQuery };
