## Installation

Ici nous installons l'application **api_gateway_sql** sous une machine linux :
- via **docker** : installation testée sur Ubuntu 20.04, Ubuntu 22.04
- via **podman** : installation testée sur Rocky linux 8.9

En prérequis, il est nécessaire d'installer ou d'utiliser un SGBD compatible, comme **MariaDB**, **MySQL**, **PostgreSQL**, **SqlServer** ou **Sqlite**. **MariaDB** est utilisé ici à titre d'exemple pour un environnement de test. Vous pouvez, par exemple, opter pour une installation conteneurisée en fonction de votre système d'exploitation. Le lien ci-dessous vous guidera pour mettre en place une sandbox **MariaDB** avec la base de données **school** :

[https://github.com/willbrid/api_gateway_sql/blob/main/fixtures/docs/databases.md](https://github.com/willbrid/api_gateway_sql/blob/main/fixtures/docs/databases.md).

A présent installons l'application **api_gateway_sql** en conteneur.

```
mkdir $HOME/api_gateway_sql && $HOME/api_gateway_sql/data && cd $HOME/api_gateway_sql
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

- **Installation sans persistence des données sqlite de l'application**

Sous Ubuntu
```
docker run -d --network=host --name api_gateway_sql -v $HOME/api_gateway_sql/config.yaml:/etc/api-gateway-sql/config.yaml -e API_GATEWAY_SQL_ENABLE_HTTPS=true willbrid/api-gateway-sql:latest
```

ou

Sous Rocky
```
podman run -d --net=host --name api_gateway_sql -v $HOME/api_gateway_sql/config.yaml:/etc/api-gateway-sql/config.yaml:z -e API_GATEWAY_SQL_ENABLE_HTTPS=true willbrid/api-gateway-sql:latest
```

- **Installation avec persistence des données sqlite de l'application**

Sous Ubuntu
```
docker run -d --network=host --name api_gateway_sql -v $HOME/api_gateway_sql/data:/data -v $HOME/api_gateway_sql/config.yaml:/etc/api-gateway-sql/config.yaml -e API_GATEWAY_SQL_ENABLE_HTTPS=true willbrid/api-gateway-sql:latest
```

ou

Sous Rocky
```
podman run -d --net=host --name api_gateway_sql -v $HOME/api_gateway_sql/data:/data:z -v $HOME/api_gateway_sql/config.yaml:/etc/api-gateway-sql/config.yaml:z -e API_GATEWAY_SQL_ENABLE_HTTPS=true willbrid/api-gateway-sql:latest
```

Une fois l'installation terminée, pour ouvrir le swagger via un navigateur, nous accédons à sa page via l'url ci-dessous

```
https://localhost:5297/swagger/index.html
```