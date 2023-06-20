class DatabaseQueryExecutor {
  static async executeQuery(dbConfig, query, params) {
    const Database = require(dbConfig.module || `@aux4/db-${dbConfig.type}`);
    const database = new Database(dbConfig);

    await database.open();
    const result = await database.execute(query, params);
    await database.close();

    return result.data;
  }
}

module.exports = DatabaseQueryExecutor;
