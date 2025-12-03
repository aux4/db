# aux4 db driver - Input Stream Tests

## execute with input stream (single JSON object)

```file:user-data.json
{"id": 1, "name": "alice", "email": "alice@example.com"}
```

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-input-stream-single",
          "execute": [
            "cat user-data.json | aux4 db driver --action execute --query \"insert into users (id, name, email) values (:id, :name, :email)\""
          ],
          "help": {
            "text": "test execute with single JSON object input stream"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-input-stream-single
```

```expect
{"action":"execute","sql":"insert into users (id, name, email) values (:id, :name, :email)","params":{},"ignore":false}
```

## execute with input stream (JSON array)

```file:users-array.json
[
  {"id": 1, "name": "bob", "email": "bob@example.com"},
  {"id": 2, "name": "carol", "email": "carol@example.com"}
]
```

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-input-stream-array",
          "execute": [
            "cat users-array.json | aux4 db driver --action execute --query \"insert into users (id, name, email) values (:id, :name, :email)\""
          ],
          "help": {
            "text": "test execute with JSON array input stream"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-input-stream-array
```

```expect
{"action":"execute","sql":"insert into users (id, name, email) values (:id, :name, :email)","params":{},"ignore":false}
```

## stream with input stream and transaction

```file:log-entries.json
[
  {"level": "info", "message": "User logged in", "timestamp": "2024-01-01T10:00:00Z"},
  {"level": "error", "message": "Failed to connect", "timestamp": "2024-01-01T10:01:00Z"},
  {"level": "info", "message": "User logged out", "timestamp": "2024-01-01T10:02:00Z"}
]
```

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-stream-input-tx",
          "execute": [
            "cat log-entries.json | aux4 db driver --action stream --query \"insert into logs (level, message, timestamp) values (:level, :message, :timestamp)\""
          ],
          "help": {
            "text": "test stream with input stream and transaction"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-stream-input-tx
```

```expect
{"action":"stream","sql":"insert into logs (level, message, timestamp) values (:level, :message, :timestamp)","params":{},"ignore":false}
```

## execute with input stream, base parameters, and ignore

```file:orders.json
[
  {"order_id": 101, "amount": 99.99},
  {"order_id": 102, "amount": 149.50}
]
```

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-input-stream-base-params",
          "execute": [
            "cat orders.json | aux4 db driver --action execute --query \"insert into orders (id, amount, user_id, status) values (:order_id, :amount, :user_id, :status)\" --params value(*) --ignore true"
          ],
          "help": {
            "text": "test execute with input stream, base parameters, and ignore"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-input-stream-base-params --user_id 123 --status pending
```

```expect
{"action":"execute","sql":"insert into orders (id, amount, user_id, status) values (:order_id, :amount, :user_id, :status)","params":{"status":"pending","user_id":"123"},"ignore":true}
```