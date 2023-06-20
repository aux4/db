const { ConfigLoader } = require("@aux4/config");

async function createDatabaseConfig(params) {
  const configFile = await params.configFile;
  const configName = await params.config;

  if (configName || configFile) {
    const config = ConfigLoader.load(configFile);
    return config.get(configName);
  }

  return {
    module: await params.module,
    type: await params.type,
    host: await params.host,
    port: await params.port,
    user: await params.user,
    password: await params.password,
    database: await params.database
  };
}

async function getQueryParams(dbConfig, params) {
  const queryParams = {};
  for (const key in params) {
    if (dbConfig[key] === undefined) {
      queryParams[key] = await params[key];
    }
  }
  return queryParams;
}

module.exports = { createDatabaseConfig, getQueryParams };
