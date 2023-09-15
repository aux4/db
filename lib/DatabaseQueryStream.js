const { getQuery } = require("./QueryLoader");
const fs = require("fs");

class DatabaseQueryStream {
  static async stream(dbConfig, query, params) {
    const Database = require(dbConfig.module || `@aux4/db-${dbConfig.type}`);
    const database = new Database(dbConfig);

    await database.open();
    const stream = await database.stream(query, params);
    stream.on("data", row => console.log(JSON.stringify(row, null, 2)));
    stream.on("error", err => {
      if (params.outputError) {
        fs.appendFileSync(params.outputError, err.message, { encoding: "utf8" });
      } else {
        console.error(err.message.red);
      }
    });
    stream.on("close", async () => await database.close());
    return stream;
  }
}

module.exports = DatabaseQueryStream;
