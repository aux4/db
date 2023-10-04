const { getQuery } = require("../../lib/QueryLoader");
const { DatabaseQueryStream, DatabaseStreamTransformer } = require("../../index");
const { readJsonInputStream } = require("../util/InputStreamUtils");
const { createDatabaseConfig, getQueryParams } = require("../util/DatabaseUtils");

async function streamExecutor(params) {
  const dbConfig = await createDatabaseConfig(params);
  const queryParams = await getQueryParams(dbConfig, params.params);

  const file = await params.file;
  const query = await params.query;
  const sql = await getQuery(file, query);
  const inputStreamParam = await params.inputStream;
  const inputStream = inputStreamParam === true || inputStreamParam === "true";

  if (inputStream) {
    await DatabaseQueryStream.stream(dbConfig, sql, queryParams);
  } else {
    const executeQueryStream = new DatabaseStreamTransformer(dbConfig, sql, queryParams, true);
    await executeQueryStream.open();

    readJsonInputStream().pipe(executeQueryStream).pipe(process.stdout);
  }
}

module.exports = { streamExecutor };
