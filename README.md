# @aux4/db

## Install

```bash
$ npm install --global @aux4/db
...
$ npm install --global @aux4/db-oracle
$ npm install --global @aux4/db-mssql
$ npm install --global @aux4/db-postgres
```

## Usage

### Execute Query

```bash
$ db execute --host localhost --port 1521 --user sysadmin --service ORCL --query "select * from schema.table where id = :id" --id 1
```

### Stream Query

```bash
$ db stream --host localhost --port 1521 --user sysadmin --service ORCL --query "select * from schema.table where id = :id" --id 1
```

### Using @aux4/config

create `config.yaml`

```yaml
config:
  dev:
    oracle:
      host: localhost
      port: 1521
      user: sysadmin
      password: "******"
      service: ORCL
```

```bash
$ db execute --configFile config.yaml --config dev/oracle --query "select * from schema.table where id = :id" --id 1
```

## See Also

* aux4 [website](https://aux4.io) / [npm](https://www.npmjs.com/package/aux4) / [GitHub](https://github.com/aux4/aux4)
* @aux4/config [npm](https://www.npmjs.com/package/@aux4/config) / [GitHub](https://github.com/aux4/config)