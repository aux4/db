const fs = require("fs");

class DatabaseQueryExecutor {
  static async executeQuery(dbConfig, query, params) {
    const Database = require(dbConfig.module || `@aux4/db-${dbConfig.type}`);
    const database = new Database(dbConfig);

    await database.open();

    try {
      const result = await database.execute(query, params);
      return result.data;
    } catch (err) {
      const outputError = await params.outputError;
      if (outputError) {
        fs.appendFileSync(outputError, err.message, { encoding: "utf8" });
      } else {
        console.error(err.message.red);
      }
    } finally {
      await database.close();
    }
  }
}

module.exports = DatabaseQueryExecutor;
