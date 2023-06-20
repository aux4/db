#!/usr/bin/env node

const { Engine } = require("@aux4/engine");
const { queryExecutor } = require("./command/QueryExecutor");
const { streamExecutor } = require("./command/StreamExecutor");

const QUERY_HELP = {
  variables: [
    {
      name: "configFile",
      text: "Configuration file"
    },
    {
      name: "config",
      text: "Configuration name"
    },
    {
      name: "module",
      text: "Database module"
    },
    {
      name: "type",
      text: "Database vendor"
    },
    {
      name: "host",
      text: "Database host"
    },
    {
      name: "port",
      text: "Database port"
    },
    {
      name: "user",
      text: "Database user",
      hide: true
    },
    {
      name: "password",
      text: "Database password",
      hide: true
    },
    {
      name: "database",
      text: "Database name"
    },
    {
      name: "query",
      text: "SQL query to execute"
    }
  ]
};

const config = {
  profiles: [
    {
      name: "main",
      commands: [
        {
          name: "execute",
          execute: queryExecutor,
          help: {
            ...QUERY_HELP,
            text: "Execute a SQL query"
          }
        },
        {
          name: "stream",
          execute: streamExecutor,
          help: {
            ...QUERY_HELP,
            text: "Execute a SQL query"
          }
        }
      ]
    }
  ]
};

(async () => {
  const engine = new Engine({ aux4: config });

  const args = process.argv.splice(2);
  await engine.run(args);
})();
