# db

This is a database scaffold for creating custom database command-line interfaces (CLIs).

## Databases

- [PostgreSQL](/r/public/packages/aux4/db-postgresql)
- [MSSQL](/r/public/packages/aux4/db-mssql)
- [Oracle](/r/public/packages/aux4/db-oracle)
- [MySQL](/r/public/packages/aux4/db-mysql)
- [SQLite/Turso](/r/public/packages/aux4/db-sqlite)
- [Other databases](/r/public/search?q=tag:database)

## DB Driver Implementation Requirements

When implementing a database driver for the aux4 db system, your driver must handle two primary operations: `execute` and `stream`. Both operations receive standardized input and must produce specific output formats.

### Input Format

The db driver receives input via stdin as JSON with the following structure:

```json
{
  "action": "execute|stream",
  "query": "SQL query string",
  "file": "path/to/sql/file (optional)",
  "inputStream": "true|false",
  "tx": "true|false", 
  "ignore": "true|false",
  "params": {
    "param1": "value1",
    "param2": "value2"
  }
}
```

When `inputStream` is `true`, the driver should also read JSON objects from stdin for batch processing:
```json
{"param1": "value1", "param2": "value2"}
{"param1": "value3", "param2": "value4"}
```

### Execute Operation

The `execute` operation should return complete result sets as JSON arrays.

**Success Output:**
```json
[
  {"column1": "value1", "column2": "value2"},
  {"column1": "value3", "column2": "value4"}
]
```

**Error Output (to stderr):**
```json
[{"item": {...}, "query": "SQL query", "error": "error message"}]
```

**With `ignore` flag:** Continue processing remaining items, output successful results to stdout, errors to stderr.

**With `tx` flag:** Execute all operations within a single transaction. On error, rollback all changes.

### Stream Operation

The `stream` operation should return results as newline-delimited JSON objects (NDJSON).

**Success Output:**
```json
{"column1": "value1", "column2": "value2"}
{"column1": "value3", "column2": "value4"}
```

**Error Output (to stderr):**
```json
{"item": {...}, "query": "SQL query", "error": "error message"}
```

**With `ignore` flag:** Continue streaming remaining results, output errors to stderr.

**With `tx` flag:** Execute within transaction, rollback on any error.

### Parameter Handling

- Parameters in queries use named placeholders (e.g., `:paramName`, `$paramName`, or `?` depending on database)
- CLI parameters override JSON input parameters
- Missing parameters should be treated as `null`
- Support both single-record and batch operations via `inputStream`

### Transaction Behavior

When `tx` is `true`:
- Begin transaction before first operation
- Commit on successful completion of all operations
- Rollback on any error
- For streaming: buffer results until transaction completes

### Error Handling

When `ignore` is `false` (default):
- Stop processing on first error
- Return error details to stderr
- Exit with non-zero code

When `ignore` is `true`:
- Continue processing remaining items
- Output successful results normally
- Send errors to stderr but continue
- Exit with zero code if any operations succeeded

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
          "execute": [
            "stdin:aux4 db driver stream params(query, file, inputStream, tx, ignore) --params value(*) | <command to execute a query> values(<db params>)"
          ],
          "help": {
            "text": "Execute a query on the <database name> database",
            "variables": [
              {
                <db variables>
              },
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
              },
              {
                "name": "tx",
                "text": "Execute the query in a transaction",
                "default": "false"
              },
              {
                "name": "ignore",
                "text": "Ignore errors and continue processing",
                "default": "false"
              }
            ]
          }
        },
        {
          "name": "stream",
          "execute": [
            "stdin:aux4 db driver stream params(query, file, inputStream, tx, ignore) --params value(*) | <command to stream a query> values(<db params>)"
          ],
          "help": {
            "text": "Execute a query on the <database name> database and stream results",
            "variables": [
              {
                <db variables>
              },
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
              },
              {
                "name": "tx",
                "text": "Execute the query in a transaction",
                "default": "false"
              },
              {
                "name": "ignore",
                "text": "Ignore errors and continue processing",
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
        "text": "Database host",
        "default": "localhost"
      },
      {
        "name": "port",
        "text": "Database port",
        "default": "<default port>"
      },
      {
        "name": "user",
        "text": "Database user",
        "default": "<default user>"
      },
      {
        "name": "password",
        "text": "Database password"
      },
      {
        "name": "database",
        "text": "Database name",
        "default": "<default database>"
      }
...
```
