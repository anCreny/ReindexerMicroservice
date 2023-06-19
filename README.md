# Reindexer Microservice
This microservice was developed to prodivde CRUD functionality for Reindexer database.

## Install
There's simple way how to install it on your device:
1. Clone repository
2. Write in console **(\*sudo) make start-database** and **(\*sudo) make start-microservice**

Req:
- Docker/docker-compose
- Git
- make
- Golang

\* On linux: if u didn't add docker/docker-compose package in root, you should use **sudo** before make to can use commands

## Description

### CRUD
After start the microservice you can use the following path to opete CRUD:

- **CREATE** to create new document in database u should use `http://localhost:your_port*/createdocument`
- **READ** to read documents from database u can use 2 ways:
  - `http://localhost:your_port*/readdocuments`: will be read all documents from database in JSON format
  - `http://localhost:your_port*/readonedocument?id=`: will be read only one document with appropriate id (will be responded HTTP 404 if id was incorrect) in JSON format
- **UPDATE** to update document in database u should use `http://localhost:your_port*/updatedocument`
- **DELETE** to delete document in database u can use `http://localhost:your_port*/deletedocument?id=` (will be responded HTTP 400 if id was incorrect)

\* see **Configs**

### Others
- If you read one document, it would be cached with timeout 15 minutes. Methods **DELETE** and **UPDATE** touch apon cache in real time
- **READ** method has pagination, you can use query params to provide it: `http://localhost:your_port/readdocuments?limit=5&page=2`:
  - *page* - number of page
  - *limit* - how many documents will be on each page



## Configs
The microservice supports config. You can configure the following options:

- **port** (port inside docker container for listening inside server, use -p flag with *port* to provide port that you want)
- **database_connection_port** (port for connection to database)
- **database_connection_username** (username for connection to database)
- **database_connection_password** (password for connection to database)
- **database_connection_name** (name of database you want to connect)

You can prodive config in two ways:
1. Via `config.yaml` file:
  Edit `config.yaml` file to prodive configs
  Fields:
  ```YAML
  port: 8080
  db_port: 6543
  db_username: userr
  db_password: pass
  db_name: db
  ```
  Path to file: `root/build/pkg/configs/config.yaml`

2. Via environment variables:
  Edit `docker-compose.yml` file to prodive configs
  Fields:
  ```docker-compose.yml
  environment:
      MICROSERVICE_PORT: 8080
      DB_CONNECTION_PORT: 6534
      DB_CONNECTION_USERNAME: user
      DB_CONNECTION_PASSWORD: pass
      DB_CONNECTION_NAME: db
  ```
  Path to file: `root/build/pkg/docker-compose.yml`

### P.S.
1. Default values for config variable you can see in both env and yaml examples

2. Configs follow the following action pattern:
  If env variable empty -> check variable from YAML file, if it also empty -> use default value for the variable
