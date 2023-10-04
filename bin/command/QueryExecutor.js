const { getQuery } = require("../../lib/QueryLoader");
const { readJsonInputStream } = require("../util/InputStreamUtils");
const { createDatabaseConfig, getQueryParams } = require("../util/DatabaseUtils");
const { DatabaseQueryExecutor, DatabaseStreamTransformer } = require("../..");

async function queryExecutor(params) {
  const dbConfig = await createDatabaseConfig(params);
  const queryParams = await getQueryParams(dbConfig, params.params);

  const file = await params.file;
  const query = await params.query;
  const sql = await getQuery(file, query);
  const inputStreamParam = await params.inputStream;
  const inputStream = inputStreamParam === true || inputStreamParam === "true";

  if (inputStream) {
    const result = await DatabaseQueryExecutor.executeQuery(dbConfig, sql, queryParams);
    console.log(JSON.stringify(result, null, 2));
  } else {
    const executeQueryStream = new DatabaseStreamTransformer(dbConfig, sql, queryParams);
    await executeQueryStream.open();

    readJsonInputStream().pipe(executeQueryStream).pipe(process.stdout);
  }
}

module.exports = { queryExecutor };
