# aux4 db driver prepare - Additional Parameters Tests

## execute with additional key-value parameters

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-additional-params",
          "execute": [
            "aux4 db driver prepare --action execute --query \"select * from users where name = :name and age = :age and active = :active\" --params value(*)"
          ],
          "help": {
            "text": "test execute with additional key-value parameters"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-additional-params --name john --age 25 --active true --debug false
```

```expect
{"action":"execute","sql":"select * from users where name = :name and age = :age and active = :active","params":{"active":"true","age":"25","debug":"false","name":"john"},"ignore":false}
```

## execute with mixed parameter types

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-mixed-params",
          "execute": [
            "aux4 db driver prepare --action execute --query \"select * from products where price \\u003e :minPrice and category = :category and featured = :featured\" --params value(*)"
          ],
          "help": {
            "text": "test execute with mixed parameter types"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-mixed-params --category electronics --minPrice 100.5 --featured true --count 10
```

```expect
{"action":"execute","sql":"select * from products where price u003e :minPrice and category = :category and featured = :featured","params":{"category":"electronics","count":"10","featured":"true","minPrice":"100.5"},"ignore":false}
```

## stream with additional parameters and ignore

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-stream-additional",
          "execute": [
            "aux4 db driver prepare --action stream --query \"select * from logs where level = :level and timestamp \\u003e :since\" --params value(*) --ignore true"
          ],
          "help": {
            "text": "test stream with additional parameters and ignore"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-stream-additional --level error --since "2024-01-01" --limit 100
```

```expect
{"action":"stream","sql":"select * from logs where level = :level and timestamp u003e :since","params":{"level":"error","limit":"100","since":"2024-01-01"},"ignore":true}
```