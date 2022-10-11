# REST API service "To help the traveler"

The REST API was written for the purpose of practicing programming in the Go language ([task](TASK.md))

The following concepts and technologies were used in the project:

* The Clean Architecture approach in building the structure of an application. Dependency injection technique.

* Working with the framework <a href="https://github.com/gin-gonic/gin">gin-gonic/gin</a>.

* Application configuration using the <a href="https://github.com/ilyakaznacheev/cleanenv">ilyakaznacheev/cleanenv</a> library.

* Working with environment variables using the <a href="https://github.com/joho/godotenv">joho/godotenv</a> library.

* Postgresql database and <a href="https://github.com/jmoiron/sqlx">jmoiron/sqlx</a> library for working with it.

* Unit testing (<a href="github.com/golang/mock/gomock">gomock</a>, <a href="https://github.com/zhashkevych/go-sqlxmock">go-sqlxmock</a>).

* Description of the API with swagger (<a href="https://github.com/swaggo/swag">swaggo/swag</a>)

* Running a project with Docker (docker-compose).

* Graceful Shutdown.

## Requirements
```
docker & docker-compose
```

## Run Project
```
Use `docker-compose up -d` to build and run docker containers with application and postgres-db instance
```