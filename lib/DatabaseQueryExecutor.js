const fs = require("fs");

class DatabaseQueryExecutor {
  static async executeQuery(dbConfig, query, params) {
    let Database;

    try {
      Database = require(dbConfig.module || `@aux4/db-${dbConfig.type}`);
    } catch (e) {
      throw new Error(`Database module not found: ${dbConfig.module || `@aux4/db-${dbConfig.type}`}`);
    }

    const database = new Database(dbConfig);

    await database.open();

    try {
      const result = await database.execute(query, params);
      return result.data;
    } catch (err) {
      if (params.outputError) {
        fs.appendFileSync(params.outputError, err.message, { encoding: "utf8" });
      } else {
        console.error(err.message.red);
      }
    } finally {
      await database.close();
    }
  }
}

module.exports = DatabaseQueryExecutor;
