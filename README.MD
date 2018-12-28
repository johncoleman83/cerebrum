# Cerebrum

[![Build Status](https://travis-ci.org/johncoleman83/cerebrum.svg?branch=master)](https://travis-ci.org/johncoleman83/cerebrum)
[![codecov](https://codecov.io/gh/johncoleman83/cerebrum/branch/master/graph/badge.svg)](https://codecov.io/gh/johncoleman83/cerebrum)
[![Go Report Card](https://goreportcard.com/badge/github.com/johncoleman83/cerebrum)](https://goreportcard.com/report/github.com/johncoleman83/cerebrum)
[![Maintainability](https://api.codeclimate.com/v1/badges/3e7ca0599c269e5ebfb6/maintainability)](https://codeclimate.com/github/johncoleman83/cerebrum/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/3e7ca0599c269e5ebfb6/test_coverage)](https://codeclimate.com/github/johncoleman83/cerebrum/test_coverage)

## Branched off of GORSK - GO(lang) Restful Starter Kit

**[App started with GORSK](https://www.ribice.ba/refactoring-gorsk/)**

Gorsk is a Golang starter kit for developing RESTful services. It is designed to help you kickstart your project, skipping the 'setting-up part' and jumping straight to writing business logic.

Previously Gorsk was built using [Gin](https://github.com/gin-gonic/gin). Gorsk using Gin is available [HERE](https://github.com/ribice/gorsk-gin).


The GORSK starter kit provided:

* Fully featured RESTful endpoints for authentication, changing password and CRUD operations on the user entity
* JWT authentication and session
* Application configuration via config file (yaml)
* RBAC (role-based access control)
* Structured logging
* Great performance
* Partial update (PATCH) handling using reflection
* Request marshaling and data validation
* API Docs using SwaggerUI
* Mocking using stdlib
* Complete test coverage
* Containerized database query tests

The following dependencies are used in this project (generated using [Glice](https://github.com/ribice/glice)):

```bash
|-------------------------------------|--------------------------------------------|--------------|
|             DEPENDENCY              |                  REPOURL                   |   LICENSE    |
|-------------------------------------|--------------------------------------------|--------------|
| github.com/labstack/echo            | https://github.com/labstack/echo           | MIT          |
| github.com/jinzhu/gorm              | https://github.com/jinzhu/gorm             | MIT          |
| github.com/dgrijalva/jwt-go         | https://github.com/dgrijalva/jwt-go        | MIT          |
| github.com/rs/zerolog               | https://github.com/rs/zerolog              | MIT          |
| golang.org/x/crypto/bcrypt          | https://github.com/golang/crypto           |              |
| gopkg.in/yaml.v2                    | https://github.com/go-yaml/yaml            |              |
| gopkg.in/go-playground/validator.v8 | https://github.com/go-playground/validator | MIT          |
| github.com/nbutton23/zxcvbn-go      | https://github.com/nbutton23/zxcvbn-go     | MIT          |
| github.com/ory/dockertest           | https://github.com/ory/dockertest          | Apache       |
| github.com/stretchr/testify         | https://github.com/stretchr/testify        | Other        |
|-------------------------------------|--------------------------------------------|--------------|
```

1. Echo - HTTP 'framework'.
2. GORM - Golang ORM
3. lib/mysql - mysql driver (part of GORM)
4. JWT-Go - JWT Authentication
5. Zerolog - Structured logging
6. Bcrypt - Password hashing
7. Yaml - Unmarshalling YAML config file
8. Validator - Request validation.
9. zxcvbn-go - Password strength checker
10. DockerTest - Testing database queries
11. Testify/Assert - Asserting test results

Most of these can easily be replaced with your own choices since their usage is abstracted and localized.

## Getting started

Using Gorsk requires having Go 1.7 or above. Once you downloaded Gorsk (either using Git or go get) you need to configure the following:

1. In cmd/migration/main.go set up psn variable and then run it (go run main.go). It will create all tables, and necessery data, with a new account username/password admin/admin.

* if you start a postgres docker container with:
  ```
  docker run --name some_postgress_name --env POSTGRES_PASSWORD=postgres --env POSTGRES_USER=postgres --publish 5432:5432 --detach postgres:alpine
  ```
  set the psn to:
  ```
  postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
  ```

* if you start with mysql:
  ```
  $ docker-compose up --detach
  ```

2. Run the app using:

```bash
go run cmd/api/main.go
```

3. Visit `http://localhost:8080/swaggerui/` (with trailing slash) in your browser to test out some endpoints!

The application runs as an HTTP server at port 8080. It provides the following RESTful endpoints:

* `POST /login`: accepts username/passwords and returns jwt token and refresh token
* `GET /refresh/:token`: refreshes sessions and returns jwt token
* `GET /me`: returns info about currently logged in user
* `GET /swaggerui/` (with trailing slash): launches swaggerui in browser
* `GET /v1/users`: returns list of users
* `GET /v1/users/:id`: returns single user
* `POST /v1/users`: creates a new user
* `PATCH /v1/password/:id`: changes password for a user
* `DELETE /v1/users/:id`: deletes a user

* You can log in as admin to the application by sending a post request to localhost:8080/login with username `admin` and password `admin` in JSON body.

* When sending the Authorization header use this format: `Authorization: Bearer TOKEN_HASH`

### Implementing CRUD of another table

Let's say you have a table named 'cars' that handles employee's cars. To implement CRUD on this table you need:

1. Inside `pkg/utl/model` create a new file named `car.go`. Inside put your entity (struct), and methods on the struct if you need them.

2. Create a new `car` folder in the (micro)service where your service will be located, most probably inside `api`. Inside create a file/service named car.go and test file for it (`car/car.go` and `car/car_test.go`). You can test your code without writing a single query by mocking the database logic inside /mock/mockdb folder. If you have complex queries interfering with other entities, you can create in this folder other files such as car_users.go or car_templates.go for example.

3. Inside car folder, create folders named `platform`, `transport` and `logging`.

4. Code for interacting with a platform like database (postgresql) should be placed under `car/platform/mysqldb`. (`pkg/api/car/platform/mysqldb/car.go`)

5. In `pkg/api/car/transport` create a new file named `http.go`. This is where your handlers are located. Under the same location create http_test.go to test your API.

6. In logging directory create a file named `car.go` and copy the logic from another service. This serves as request/response logging.

6. In `pkg/api/api.go` wire up all the logic, by instantiating car service, passing it to the logging and transport service afterwards.

### Implementing other platforms

Similarly to implementing APIs relying only on a database, you can implement other platforms by:

1. In the service package, in car.go add interface that corresponds to the platform, for example, Indexer or Reporter.

2. Rest of the procedure is same, except that in `/platform` you would create a new folder for your platform, for example, `elastic`.

3. Once the new platform logic is implemented, create an instance of it in main.go (for example `elastic.Client`) and pass it as an argument to car service (`pkg/api/car/car.go`).

### Running database queries in transaction

To use a transaction, before interacting with db create a new transaction:

```go
err := s.db.RunInTransaction(func (tx *pg.Tx) error{
    // Application service here
})
````

Instead of passing database client as `s.db` , inside this function pass it as `tx`. Handle the error accordingly.

## Project Structure

Follows this [project layout](https://github.com/golang-standards/project-layout)

1. Root directory contains things not related to code directly, e.g. docker-compose, CI/CD, readme, bash scripts etc. It should also contain vendor folder, Gopkg.toml and Gopkg.lock if dep is being used.

2. Cmd package contains code for starting applications (main packages). The directory name for each application should match the name of the executable you want to have. Gorsk is structured as a monolith application but can be easily restructured to contain multiple microservices. An application may produce multiple binaries, therefore Gorsk uses the Go convention of placing main package as a subdirectory of the cmd package. As an example, scheduler applications binary would be located under cmd/cron. It also loads the necessery configuration and passes it to the service initializers.

3. Rest of the code is located under /pkg. The pkg directory contains `utl` and 'microservice' directories.

4. Microservice directories, like api (naming corresponds to `cmd/` folder naming) contains multiple folders for each domain it interacts with, for example: user, car, appointment etc.

5. Domain directories, like user, contain all application/business logic and two additional directories: platform and transport.

6. Platform folder contains various packages that provide support for things like databases, authentication or even marshaling. Most of the packages located under platform are decoupled by using interfaces. Every platform has its own package, for example, postgres, elastic, redis, memcache etc.

7. Transport package contains HTTP handlers. The package receives the requests, marshals, validates then passes it to the corresponding service.

8. Utl directory contains helper packages and models. Packages such as mock, middleware, configuration, server are located here.

## Make

* use make to help with development

```
make refresh
```

* clean up stuff
```
make clean_dev
make clean_test
```


## Tests

* trouble with DB connection, then add this to docker-compose and remove synchronous db testing:
```
  db_test:
    image: mysql:latest
    container_name: cerebrum_mysql_test_db_no_0
    command: --default-authentication-plugin=mysql_native_password
    restart: always
    environment:
      - MYSQL_ALLOW_EMPTY_PASSWORD=yes
      - MYSQL_DATABASE=cerebrum_mysql_test_db
      - MYSQL_USER=mysql_test_user
      - MYSQL_PASSWORD=mysql_test_password
    ports:
      - 3307:3306
```


## License

gorsk is licensed under the MIT license. Check the [LICENSE](LICENSE.md) file for details.

## Author

[David John Coleman II](https://davidjohncoleman.com)
