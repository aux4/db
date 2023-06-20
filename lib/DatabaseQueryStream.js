const { getQuery } = require("./QueryLoader");

class DatabaseQueryStream {
  static async stream(dbConfig, query, params) {
    const Database = require(dbConfig.module || `@aux4/db-${dbConfig.type}`);
    const database = new Database(dbConfig);

    await database.open();
    const stream = await database.stream(query, params);
    stream.on("data", row => console.log(JSON.stringify(row, null, 2)));
    stream.on("error", err => console.error(err.message));
    stream.on("close", async () => await database.close());
    return stream;
  }
}

module.exports = DatabaseQueryStream;
