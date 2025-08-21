# db

This is a database scaffold for creating custom database command-line interfaces (CLIs).

## Databases

- [PostgreSQL](/r/public/packages/aux4/db-postgresql)
- [MSSQL](/r/public/packages/aux4/db-mssql)
- [Oracle](/r/public/packages/aux4/db-oracle)
- [MySQL](/r/public/packages/aux4/db-mysql)
- [SQLite/Turso](/r/public/packages/aux4/db-sqlite)
- [Other databases](/r/public/search?q=tag:database)

## How to Create a Custom Database CLI

```json
{
  "scope": "<scope>",
  "name": "db-<database name>",
  "version": "<version>",
  "description": "<database name> database tools",
  "tags": ["db", "database", "<database name>"],
  "dependencies": ["aux4/db"],
  "profiles": [
    {
      "name": "db",
      "commands": [
        {
          "name": "<database name>",
          "execute": ["profile:db:<database name>"],
          "help": {
            "text": "<database name> database tools"
          }
        }
      ]
    },
    {
      "name": "db:<database name>",
      "commands": [
        {
          "name": "execute",
          "execute": ["<command to execute a query>"],
          "help": {
            "text": "Execute a query on the <database name> database",
            "variables": [
              {
                "name": "query",
                "text": "SQL query to execute",
                "default": ""
              },
              {
                "name": "file",
                "text": "SQL file to execute",
                "default": ""
              },
              {
                "name": "inputStream",
                "text": "Read JSON input from stdin",
                "default": "false"
              }
            ]
          }
        },
        {
          "name": "stream",
          "execute": ["<command to execute a query>"],
          "help": {
            "text": "Execute a query on the <database name> database and stream results",
            "variables": [
              {
                "name": "query",
                "text": "SQL query to execute",
                "default": ""
              },
              {
                "name": "file",
                "text": "SQL file to execute",
                "default": ""
              },
              {
                "name": "inputStream",
                "text": "Read JSON input from stdin",
                "default": "false"
              }
            ]
          }
        }
      ]
    }
  ]
}
```

You can also include the variables required to connect to the database in the `execute` and `stream` commands, for example:

```json
...
      {
        "name": "host",
        "text": "Database host"
      },
      {
        "name": "port",
        "text": "Database port"
      },
      {
        "name": "user",
        "text": "Database user"
      },
      {
        "name": "password",
        "text": "Database password"
      },
      {
        "name": "database",
        "text": "Database name"
      }
...
```
