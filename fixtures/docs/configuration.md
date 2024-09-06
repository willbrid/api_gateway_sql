## Fichier de configuration

```
api_gateway_sql:
  # Configuration de la base de données
  sqlitedb: "api_gateway_sql"
  # Configuration des paramètres d'authentification
  auth:
    # Paramètre d'activation ou de désactivation de l'authentification
    enabled: true
    # Paramètre d'utilisateur pris en compte lorsque l'authentification est activée
    username: test
    # Paramètre de mot de passe pris en compte lorsque l'authentification est activée
    password: test@test
  # Configuration des paramètres de base de données cible
  databases:
    # Paramètre d'identifiant de la cible
  - name: school
    # Paramètre du type de sgbd
    type: mariadb
    # Paramètre d'adresse de la base de données
    host: "@HOST_IP"
    # Paramètre de port de la base de données
    port: 3307
    # Paramètre d'utilisateur de la base de données
    username: "test"
    # Paramètre de mot de passe d'utilisateur de la base de données
    password: "test"
    # Paramètre du nom de la base de données
    dbname: "school"
    # Paramètre d'activation ou de désactivation du mode ssl de communication avec la base de données
    sslmode: false
    # Paramètre de timeout de communication avec la base de données
    timeout: 1s
  # Configuration des paramètres des cibles
  targets:
    # Paramètre de nom de la cible
  - name: insert_batch_student
    # Paramètre de nom pour la base de données cible
    data_source_name: school
    # Paramètre d'activation ou de désactivation de l'exécution en masse
    multi: true
    # Paramètre de la taille d'un batch à exécuter. Utiliser lorsque l'exécution en masse est activée
    batch_size: 10
    # Paramètre du nombre de blocs à utiliser pour décomposer le fichier csv. Utiliser lorsque l'exécution en masse est activée
    buffer_size: 50
    # Paramètre de champs d'une table en base de données. Utiliser lorsque l'exécution en masse est activée
    batch_fields: "name;address"
    # Paramètre de contenu d'une requête sql
    sql: "insert into school (name, address) values ({{name}}, {{address}})"
```