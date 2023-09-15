#!/usr/bin/env node

const { Engine } = require("@aux4/engine");
const { queryExecutor } = require("./command/QueryExecutor");
const { streamExecutor } = require("./command/StreamExecutor");

const QUERY_HELP = {
  variables: [
    {
      name: "configFile",
      text: "Configuration file",
      default: ""
    },
    {
      name: "config",
      text: "Configuration name",
      default: ""
    },
    {
      name: "module",
      text: "Database module",
      default: ""
    },
    {
      name: "type",
      text: "Database vendor",
      default: ""
    },
    {
      name: "host",
      text: "Database host",
      default: ""
    },
    {
      name: "port",
      text: "Database port",
      default: ""
    },
    {
      name: "user",
      text: "Database user",
      default: "",
      hide: true
    },
    {
      name: "password",
      text: "Database password",
      default: "",
      hide: true
    },
    {
      name: "database",
      text: "Database name",
      default: ""
    },
    {
      name: "query",
      text: "SQL query to execute",
      default: ""
    },
    {
      name: "file",
      text: "SQL file to execute",
      default: ""
    },
    {
      name: "errorOutput",
      text: "Path of error output file. If not provided it will output to stderr",
      default: ""
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

  try {
    await engine.run(args);
  } catch (e) {
    console.error(e.message.red);
    process.exit(1);
  }
})();
