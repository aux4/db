# aux4 db driver sql-bind - Advanced SQL Parsing Tests

## complex query with multiple parameters

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-sql-bind-complex",
          "execute": [
            "aux4 db driver sql-bind --query \"WITH recent_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders WHERE created_at > :since GROUP BY user_id) SELECT u.name, u.email, ro.order_count FROM users u JOIN recent_orders ro ON u.id = ro.user_id WHERE ro.order_count > :min_orders AND u.status = :status\" --params '{\"since\": \"2024-06-01\", \"min_orders\": 5, \"status\": \"active\"}'"
          ],
          "help": {
            "text": "test sql-bind with complex CTE query"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-sql-bind-complex
```

```expect
WITH recent_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders WHERE created_at > '2024-06-01' GROUP BY user_id) SELECT u.name, u.email, ro.order_count FROM users u JOIN recent_orders ro ON u.id = ro.user_id WHERE ro.order_count > 5 AND u.status = 'active'
```

## array with mixed types

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-sql-bind-mixed-array",
          "execute": [
            "aux4 db driver sql-bind --query \"select * from products where (id IN :ids OR price IN :prices) and active = :active\" --params '{\"ids\": [1, 2, 3], \"prices\": [9.99, 19.99, 29.99], \"active\": true}'"
          ],
          "help": {
            "text": "test sql-bind with mixed array types"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-sql-bind-mixed-array
```

```expect
select * from products where (id IN (1, 2, 3) OR price IN (9.99, 19.99, 29.99)) and active = TRUE
```

## multiple occurrences of same parameter

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-sql-bind-repeated-param",
          "execute": [
            "aux4 db driver sql-bind --query \"select * from logs where (level = :level OR fallback_level = :level) and message like :pattern and created_at > :date and updated_at > :date\" --params '{\"level\": \"error\", \"pattern\": \"%failed%\", \"date\": \"2024-01-01\"}'"
          ],
          "help": {
            "text": "test sql-bind with repeated parameters"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-sql-bind-repeated-param
```

```expect
select * from logs where (level = 'error' OR fallback_level = 'error') and message like '%failed%' and created_at > '2024-01-01' and updated_at > '2024-01-01'
```

## array with simple strings

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-sql-bind-array-strings",
          "execute": [
            "aux4 db driver sql-bind --query \"select * from comments where content IN :messages\" --params '{\"messages\": [\"Hello world\", \"Good morning\", \"Have a nice day\"]}'"
          ],
          "help": {
            "text": "test sql-bind with array of strings"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-sql-bind-array-strings
```

```expect
select * from comments where content IN ('Hello world', 'Good morning', 'Have a nice day')
```

## false and zero values

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-sql-bind-falsy-values",
          "execute": [
            "aux4 db driver sql-bind --query \"select * from settings where enabled = :enabled and retry_count = :retries and timeout = :timeout\" --params '{\"enabled\": false, \"retries\": 0, \"timeout\": 0.0}'"
          ],
          "help": {
            "text": "test sql-bind with false and zero values"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-sql-bind-falsy-values
```

```expect
select * from settings where enabled = FALSE and retry_count = 0 and timeout = 0
```