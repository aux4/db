# aux4 db driver - File Tests

## execute with SQL file

```file:query.sql
select name, email from users
where created_at > :startDate
  and status = :status
order by created_at desc
```

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-file",
          "execute": [
            "aux4 db driver --action execute --file query.sql --params value(*)"
          ],
          "help": {
            "text": "test execute with SQL file"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-file --startDate "2024-01-01" --status active
```

```expect
{"action":"execute","sql":"select name, email from users\nwhere created_at \u003e :startDate\n  and status = :status\norder by created_at desc","params":{"startDate":"2024-01-01","status":"active"},"ignore":false}
```

## stream with SQL file and ignore

```file:complex-query.sql
with recent_orders as (
  select user_id, count(*) as order_count
  from orders
  where created_at > :since
  group by user_id
)
select u.name, u.email, ro.order_count
from users u
join recent_orders ro on u.id = ro.user_id
where ro.order_count > :minOrders
```

```file:.aux4
{
  "profiles": [
    {
      "name": "main",
      "commands": [
        {
          "name": "test-file-stream",
          "execute": [
            "aux4 db driver --action stream --file complex-query.sql --params value(*) --ignore true"
          ],
          "help": {
            "text": "test stream with SQL file and ignore"
          }
        }
      ]
    }
  ]
}
```

```execute
aux4 test-file-stream --since "2024-06-01" --minOrders 5
```

```expect
{"action":"stream","sql":"with recent_orders as (\n  select user_id, count(*) as order_count\n  from orders\n  where created_at \u003e :since\n  group by user_id\n)\nselect u.name, u.email, ro.order_count\nfrom users u\njoin recent_orders ro on u.id = ro.user_id\nwhere ro.order_count \u003e :minOrders","params":{"minOrders":"5","since":"2024-06-01"},"ignore":true}
```