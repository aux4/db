# aux4 db driver sql-bind - Basic SQL Parsing Tests

## simple parameter substitution

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-sql-bind-simple",
          "execute": [
            "aux4 db driver sql-bind --query \"select * from users where name = :name and age = :age\" --params '{\"name\": \"alice\", \"age\": 25}'"
          ],
          "help": {
            "text": "test sql-bind with simple parameters"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-sql-bind-simple
```

```expect
select * from users where name = 'alice' and age = 25
```

## string with quotes (SQL injection protection)

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-sql-bind-injection",
          "execute": [
            "aux4 db driver sql-bind --query \"select * from users where name = :name\" --params '{\"name\": \"O'\\''Reilly\"}'"
          ],
          "help": {
            "text": "test sql-bind with quotes in data"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-sql-bind-injection
```

```expect
select * from users where name = 'O''Reilly'
```

## boolean and null values

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-sql-bind-types",
          "execute": [
            "aux4 db driver sql-bind --query \"select * from users where active = :active and manager_id = :manager_id\" --params '{\"active\": true, \"manager_id\": null}'"
          ],
          "help": {
            "text": "test sql-bind with boolean and null values"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-sql-bind-types
```

```expect
select * from users where active = TRUE and manager_id = NULL
```

## IN clause with array

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-sql-bind-in-clause",
          "execute": [
            "aux4 db driver sql-bind --query \"select * from users where id IN :ids and department IN :depts\" --params '{\"ids\": [1, 2, 3], \"depts\": [\"engineering\", \"sales\"]}'"
          ],
          "help": {
            "text": "test sql-bind with IN clause arrays"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-sql-bind-in-clause
```

```expect
select * from users where id IN (1, 2, 3) and department IN ('engineering', 'sales')
```

## empty array handling

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-sql-bind-empty-array",
          "execute": [
            "aux4 db driver sql-bind --query \"select * from users where id IN :ids\" --params '{\"ids\": []}'"
          ],
          "help": {
            "text": "test sql-bind with empty array"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-sql-bind-empty-array
```

```expect
select * from users where id IN (NULL)
```

## no parameters in query

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-sql-bind-no-params",
          "execute": [
            "aux4 db driver sql-bind --query \"select count(*) from users\" --params '{}'"
          ],
          "help": {
            "text": "test sql-bind with no parameters"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-sql-bind-no-params
```

```expect
select count(*) from users
```