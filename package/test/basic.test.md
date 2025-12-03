# aux4 db driver - Basic Tests

## basic execute with named parameters

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-execute",
          "execute": [
            "aux4 db driver --action execute --query \"select 1 from table where name = :name\" --params value(*)"
          ],
          "help": {
            "text": "test basic execute command"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-execute --name david
```

```expect
{"action":"execute","sql":"select 1 from table where name = :name","params":{"name":"david"},"ignore":false}
```

## execute with multiple parameters

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-multi-params",
          "execute": [
            "aux4 db driver --action execute --query 'select * from users where name = :name and age = :age' --params value(*)"
          ],
          "help": {
            "text": "test execute with multiple parameters"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-multi-params --name john --age 25
```

```expect
{"action":"execute","sql":"select * from users where name = :name and age = :age","params":{"age":"25","name":"john"},"ignore":false}
```

## execute with ignore flag

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-ignore",
          "execute": [
            "aux4 db driver --action execute --query \"select 1 from table where name = :name\" --params value(*) --ignore true"
          ],
          "help": {
            "text": "test execute with ignore flag"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-ignore --name david
```

```expect
{"action":"execute","sql":"select 1 from table where name = :name","params":{"name":"david"},"ignore":true}
```

## execute without parameters

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-no-params",
          "execute": [
            "aux4 db driver --action execute --query \"select 1 from dual\""
          ],
          "help": {
            "text": "test execute without parameters"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-no-params
```

```expect
{"action":"execute","sql":"select 1 from dual","params":{},"ignore":false}
```
