# Cerebrum

[![Build Status](https://travis-ci.org/johncoleman83/cerebrum.svg?branch=master)](https://travis-ci.org/johncoleman83/cerebrum)
[![codecov](https://codecov.io/gh/johncoleman83/cerebrum/branch/master/graph/badge.svg)](https://codecov.io/gh/johncoleman83/cerebrum)
[![Go Report Card](https://goreportcard.com/badge/github.com/johncoleman83/cerebrum)](https://goreportcard.com/report/github.com/johncoleman83/cerebrum)
[![Maintainability](https://api.codeclimate.com/v1/badges/3e7ca0599c269e5ebfb6/maintainability)](https://codeclimate.com/github/johncoleman83/cerebrum/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/3e7ca0599c269e5ebfb6/test_coverage)](https://codeclimate.com/github/johncoleman83/cerebrum/test_coverage)

## Branched off of GORSK - GO(lang) Restful Starter Kit

**[App started with GORSK](https://www.ribice.ba/refactoring-gorsk/)**

Gorsk is a Golang starter kit for developing RESTful services. It is designed to help you kickstart your project, skipping the 'setting-up part' and jumping straight to writing business logic.

## Dependencies

```bash
|------------------------------------|--------------|
|              PACKAGE               |    LICENSE   |
|------------------------------------|--------------|
| github.com/labstack/echo           | MIT          |
| github.com/jinzhu/gorm             | MIT          |
| github.com/dgrijalva/jwt-go        | MIT          |
| github.com/rs/zerolog              | MIT          |
| github.com/go-playground/validator | MIT          |
| github.com/nbutton23/zxcvbn-go     | MIT          |
| github.com/stretchr/testify        | MIT          |
|------------------------------------|--------------|
| github.com/go-yaml/yaml            | Apache       |
| github.com/ory/dockertest          | Apache       |
|------------------------------------|--------------|
```

## Project Structure

* Started with the [project structure](https://github.com/ribice/gorsk#project-structure) from gorsk

* trying to maintain [golang-standards](https://github.com/golang-standards/project-layout) for project layout


## Development

ensure that you have properly setup dependencies

```
$ make deps
```

use make to help with development

* just run make help to see all the automated support with `make`
  ```
  $ make help
  ```

* after setup and successfully running `make serve` use [swagger](https://github.com/johncoleman83/cerebrum#swagger) to test endpoints.

* or use `curl` to login with bootstrap data
  ```
  $ curl -X POST "http://localhost:8080/login" \
      -H "accept: application/json" \
      -H "Content-Type: application/json" \
      -d "{ \"password\": \"admin\", \"username\": \"admin\"}"
  ```

## Swagger

Visit `http://localhost:8080/swaggerui/` (with trailing slash) in your browser to test out some endpoints!

* use `make` or [multi-file-swagger](https://www.npmjs.com/package/multi-file-swagger) to generate the swagger `.yaml` spec file.
* more docs on creating a swagger server: https://goswagger.io/use/spec.html#parsing-rules
* this application uses the latest distribution from [github.com/swagger-api/swagger-ui](https://github.com/swagger-api/swagger-ui/tree/master/dist)

## Endpoints

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

## Support

### Implementing CRUD of another table

Let's say you have a table named 'cars' that handles employee's cars. To implement CRUD on this table you need:

1. Inside `pkg/utl/model` create a new file named `car.go`. Inside put your entity (struct), and methods on the struct if you need them.

2. Create a new `car` folder in the (micro)service where your service will be located, most probably inside `api`. Inside create a file/service named car.go and test file for it (`car/car.go` and `car/car_test.go`). You can test your code without writing a single query by mocking the database logic inside /mock/mockstore folder. If you have complex queries interfering with other entities, you can create in this folder other files such as car_users.go or car_templates.go for example.

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

* http://gorm.io/docs/transactions.html

## License

Currently licensed under MIT, check the [LICENSE](LICENSE.md) file for details.

## Author

[David John Coleman II](https://davidjohncoleman.com)
**[Initial Commit by GORSK](https://www.ribice.ba/refactoring-gorsk/)**
