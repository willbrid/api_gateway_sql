## Installation et utilisation

Ici nous installons l'application **api_gateway_sql** via docker sous une machine linux.

```
cd $HOME && mkdir api_gateway_sql && cd api_gateway_sql
```

```
vi config.yaml
```

```
api_gateway_sql:
  sqlitedb: "api_gateway_sql"
  auth:
    enabled: true
    username: test
    password: test@test
  databases:
  - name: school
    type: mariadb
    host: "127.0.0.1"
    port: 3307
    username: "test"
    password: "test"
    dbname: "school"
    sslmode: false
    timeout: 1s
  targets:
  - name: list-student
    data_source_name: school
    multi: false
    sql: "select * from student"
  - name: list-school
    data_source_name: school
    multi: false
    sql: "select * from school"
  - name: find-one-student
    data_source_name: school
    multi: false
    sql: "select * from student where id = {{id}}"
  - name: find-student-with-cond
    data_source_name: school
    multi: false
    sql: "select * from student where class_id = {{class}} and school_id = {{school}} and age >= {{age}}"
  - name: insert_student
    data_source_name: school
    multi: false
    sql: "insert into school (id, name, address) values ({{id}}, {{name}}, {{address}})"
  - name: insert_batch_student
    data_source_name: school
    multi: true
    batch_size: 10
    buffer_size: 50
    batch_fields: "name;address"
    sql: "insert into school (name, address) values ({{name}}, {{address}})"
```

```
docker run -d --network=host --name api_gateway_sql -v $HOME/api_gateway_sql/config.yaml:/etc/api-gateway-sql/config.yaml -e API_GATEWAY_SQL_ENABLE_HTTPS=true willbrid/api-gateway-sql:latest
```

Pour ouvrir le swagger via un navigateur, nous accédons à sa page via l'url ci-dessous

```
https://localhost:5297/swagger/index.html
```