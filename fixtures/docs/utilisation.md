## Utilisation

L'objectif de ce contenu est de fournir une présentation sur comment utiliser l'application à travers ses apis. Nous utiliserons la commande **curl** pour fournir des exemples pratiques de consommation.

#### Api [POST] : /v1/api-gateway-sql/{datasource}/init

Cette api peut être utilisé pour créer le schéma de la base de données et y insérer certaines données.

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

NB : Cette api n'est pas obligatoire si nous utilisons l'application sur des bases de données existantes.

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

#### Api [POST] : /v1/api-gateway-sql/{target}/batch

Cette API permet d'exécuter une requête SQL en mode batch en se basant sur le nom de la cible (**target**), qui contient la configuration de la requête. La requête SQL est paramétrée à l'aide d'un ou plusieurs paramètres définis par le symbole **{{}}**, et les valeurs de ces paramètres doivent être envoyées via un fichier CSV lors d'une requête POST. La configuration de la cible inclut les éléments suivants :
- Activation du mode batch (**multi: true**)
- Définition de la taille maximale d'un bloc de données du fichier CSV (**buffer_size: 50**)
- Définition de la taille maximale de chaque lot (batch) dans un bloc (**batch_size: 10**)
- Définition des champs de paramètres, où chaque champ correspond à une colonne du fichier CSV, dans l'ordre (**batch_fields: "name;address"**)

En guise de test, vous pouvez générer un fichier CSV de 100 lignes avec deux colonnes : la première colonne contient les noms des écoles, et la deuxième leur adresse. Ce fichier CSV peut être généré à l'aide d'un script Bash disponible dans ce référentiel : [https://github.com/willbrid/api_gateway_sql/blob/main/fixtures/generate_schools.sh](https://github.com/willbrid/api_gateway_sql/blob/main/fixtures/generate_schools.sh).

```
mkdir -p $HOME/api_gateway_sql && cd $HOME/api_gateway_sql
```

```
curl -fsSL https://github.com/willbrid/api_gateway_sql/raw/main/fixtures/generate_schools.sh -o generate_schools.sh
```

```
chmod +x generate_schools.sh
./generate_schools.sh 100
```

Ce script générera un fichier csv à l'emplacement **/tmp/schools.csv**.

```
curl -k -v -X POST -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -F "csvfile=@/tmp/schools.csv" https://localhost:5297/v1/api-gateway-sql/insert_batch_school/batch
```

**insert_batch_school** est le nom d'une cible configurée dans la section **api_gateway_sql.targets** du fichier de configuration. Cette cible permet d'exécuter en batch et en parallèle des insertions SQL, en récupérant les valeurs depuis le fichier **/tmp/schools.csv**.

#### Api [GET] : /v1/api-gateway-sql/stats

Cette API permet de consulter les statistiques d'exécution des requêtes en batch et de suivre leur progression.

```
curl -k -v -X GET -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' 'https://localhost:5297/v1/api-gateway-sql/stats?page_num=1&page_size=20'
```

La réponse de cette API fournit des informations sur l'exécution, notamment :
- la cible correspondante
- pour chaque bloc : son numéro de début, son numéro de fin, le nombre de succès, le nombre d'échecs, et la plage des lignes ayant échoué