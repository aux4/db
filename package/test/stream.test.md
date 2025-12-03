# aux4 db driver prepare - Stream Tests

## basic stream command

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-stream",
          "execute": [
            "aux4 db driver prepare --action stream --query \"select * from users where active = :active\" --params value(*)"
          ],
          "help": {
            "text": "test basic stream command"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-stream --active true
```

```expect
{"action":"stream","sql":"select * from users where active = :active","params":{"active":"true"},"ignore":false}
```

## stream with ignore flag

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-stream-ignore",
          "execute": [
            "aux4 db driver prepare --action stream --query \"select * from logs\" --ignore true"
          ],
          "help": {
            "text": "test stream with ignore flag"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-stream-ignore
```

```expect
{"action":"stream","sql":"select * from logs","params":{},"ignore":true}
```

## stream with complex parameters

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-stream-complex",
          "execute": [
            "aux4 db driver prepare --action stream --query \"select * from orders where user_id = :userId and status = :status and created_at > :date\" --params value(*)"
          ],
          "help": {
            "text": "test stream with complex parameters"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-stream-complex --userId 123 --status pending --date "2024-01-01"
```

```expect
{"action":"stream","sql":"select * from orders where user_id = :userId and status = :status and created_at \u003e :date","params":{"date":"2024-01-01","status":"pending","userId":"123"},"ignore":false}
```