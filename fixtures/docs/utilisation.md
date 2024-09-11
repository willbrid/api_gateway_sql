## Utilisation

L'objectif de ce contenu est de fournir une présentation sur comment utiliser l'application à travers ses apis. Nous utiliserons l'outil **postman** pour fournir des exemples pratiques de consommation.

#### Api [POST] : /v1/api-gateway-sql/{datasource}/init

Cet api peut être utilisé pour créer le schéma de la base de données et y insérer certaines données.

Pour notre environnement de test nous allons utiliser le fichier **school_mariadb.sql** qui peut être téléchargé via le lien ci-dessous : 

[https://github.com/willbrid/api_gateway_sql/blob/main/fixtures/sql/school_mariadb.sql](https://github.com/willbrid/api_gateway_sql/blob/main/fixtures/sql/school_mariadb.sql)

```
mkdir -p $HOME/api_gateway_sql && cd $HOME/api_gateway_sql
```

```
curl -fsSL https://github.com/willbrid/api_gateway_sql/raw/main/fixtures/sql/school_mariadb.sql -o school_mariadb.sql
```

```
curl -k -v -X POST -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -F "sqlfile=@school_mariadb.sql" https://localhost:5297/v1/api-gateway-sql/school/init
```

- La valeur de l'entête **Basic** représente le **base64** des crédentials (**username:password**) de l'application spécifiés au niveau de son fichier de configuration

```
echo -n test:test@test | base64
```

- **school** est le nom de la chaine de connexion à la base de données sur **Mariadb** dont nous avons configuré dans la section **api_gateway_sql.databases** du fichier de configuration de l'application.

NB : Cet api n'est pas obligatoire si nous utilisons l'application sur des bases de données existantes.

#### Api [GET] : /v1/api-gateway-sql/{target}

Cette API permet d'exécuter une requête SQL en se basant sur le nom de la cible (**target**), qui contient la configuration de la requête. Cette requête sql ne doit pas être paramétrée avec le symbole **{{}}**.

```
curl -k -v -X GET -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -H 'Content-Type: application/json' https://localhost:5297/v1/api-gateway-sql/list-student
```

```
curl -k -v -X GET -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -H 'Content-Type: application/json' https://localhost:5297/v1/api-gateway-sql/list-school
```

**list-student** et **list-school** sont les noms des cibles configurés dans la section **api_gateway_sql.targets** du fichier de configuration :
- avec **list-student** sa requête sql sélectionne toutes les lignes de la table **student**
- avec **list-school** sa requête sql sélectionne toutes les lignes de la table **school**

#### Api [POST] : /v1/api-gateway-sql/{target}

Cette API permet d'exécuter une requête SQL en se basant sur le nom de la cible (**target**), qui contient la configuration de la requête. Cette requête SQL doit être paramétrée à l'aide d'un ou plusieurs paramètres définis par le symbole **{{}}**. Les valeurs des paramètres doivent être transmises via une requête POST.

```
curl -k -v -X POST -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"id":"1"}' https://localhost:5297/v1/api-gateway-sql/find-one-student
```

```
curl -k -v -X POST -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"class":"1", "school":"1", "age":"15"}' https://localhost:5297/v1/api-gateway-sql/find-student-with-cond
```

```
curl -k -v -X POST -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"id":"300", "name":"high-tech", "address":"Willow Ave"}' https://localhost:5297/v1/api-gateway-sql/insert_school
```

**find-one-student**, **find-student-with-cond** et **insert_school** sont les noms des cibles configurés dans la section **api_gateway_sql.targets** du fichier de configuration :
- avec **find-one-student** sa requête sql sélectionne une ligne de la table **student** ayant la valeur du champ **id** égale à **1**
- avec **find-student-with-cond** sa requête sql sélectionne toutes les lignes de la table **student** ayant la valeur du champ **class_id** égale à **1**, la valeur du champ **school_id** égale à **1** et la valeur du champ **age** supérieure ou égale à **15**

NB: Chaque nom de paramètre envoyé en post doit être identique à un nom de paramètre configuré dans la requête sql.