## Démarrage des bases de données de test avec docker ou podman

#### Mariadb

```
docker run -d --name school_mariadb -p 3307:3306 --env MARIADB_USER=test --env MARIADB_PASSWORD=test --env MARIADB_DATABASE=school --env MARIADB_ROOT_PASSWORD=secret mariadb:10.5.26
```

ou

```
podman run -d --name school_mariadb -p 3307:3306 --env MARIADB_USER=test --env MARIADB_PASSWORD=test --env MARIADB_DATABASE=school --env MARIADB_ROOT_PASSWORD=secret mariadb:10.5.26
```

#### MySql

```
docker run -d --name school_mysql -p 3306:3306 --env MYSQL_USER=test --env MYSQL_PASSWORD=test --env MYSQL_DATABASE=school --env MYSQL_ROOT_PASSWORD=secret mysql:8.4.2
```

ou

```
podman run -d --name school_mysql -p 3306:3306 --env MYSQL_USER=test --env MYSQL_PASSWORD=test --env MYSQL_DATABASE=school --env MYSQL_ROOT_PASSWORD=secret mysql:8.4.2
```

#### Postgresql

```
docker run -d --name school_postgres -p 5432:5432 -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -e POSTGRES_DB=school postgres:14.3
```

ou

```
podman run -d --name school_postgres -p 5432:5432 -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -e POSTGRES_DB=school postgres:14.3
```

#### Sqlserver

```
docker run -d --name school_sqlserver --hostname school -p 1433:1433 -e "ACCEPT_EULA=Y" -e "MSSQL_SA_PASSWORD=test@Test" mcr.microsoft.com/mssql/server:2022-latest
```

ou

```
podman run -d --name school_sqlserver --hostname school -p 1433:1433 -e "ACCEPT_EULA=Y" -e "MSSQL_SA_PASSWORD=test@Test" mcr.microsoft.com/mssql/server:2022-latest
```

création de la base de données **school** sur notre instance **sqlserver**

```
docker exec -it school_sqlserver /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P test@Test -C -Q "create database school;"
```

ou 

```
podman exec -it school_sqlserver /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P test@Test -C -Q "create database school;"
```