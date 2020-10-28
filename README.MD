# Cerebrum

[![CircleCI](https://circleci.com/gh/circleci/circleci-docs.svg?style=svg)](https://circleci.com/gh/johncoleman83/cerebrum)
[![codecov](https://codecov.io/gh/johncoleman83/cerebrum/branch/main/graph/badge.svg?token=oIta2VhQc8)](https://codecov.io/gh/johncoleman83/cerebrum)
[![Go Report Card](https://goreportcard.com/badge/github.com/johncoleman83/cerebrum)](https://goreportcard.com/report/github.com/johncoleman83/cerebrum)
[![Maintainability](https://api.codeclimate.com/v1/badges/a99a88d28ad37a79dbf6/maintainability)](https://codeclimate.com/github/codeclimate/codeclimate/maintainability)

### Original design from [GORSK - GO(lang) Restful Starter Kit](https://www.ribice.ba/refactoring-gorsk/)

Gorsk is a Golang starter kit for developing RESTful services. It is designed to help you kickstart your project, skipping the 'setting-up part' and jumping straight to writing business logic.

#### Updates from GORSK:

* use `gorm` as the ORM
  * `gorm` doesn't have the limitations indicated by the author of gorsk
  * `gorm` allows for use of various different DB's
  * in this case we use **MySQL**
* get rid of go-swagger
  * it's too complicated to find Open API docs all throughout the codebase comments
  * built a new design to organize and compile all Open API doc in `third_party/swaggerui/spec`
* Updated the swagger distribution and spec file that were out of date.
* new naming conventions / refactoring
  * **Account** and **Team** models created and therefore removed **Company** and **Location**
  * Other packages and directories have been renamed to help ease the understanding of what they are used for
  * other minor refactoring and renaming of directories such as moving the swagger distribution from `/assets` to `/third_party/swaggerui/dist`
* Use CirclCI instead of Travis
* There's no reason to start 10 new docker containers for every test DB test.
  * Just use 1 test container and be sure to refresh the DB for every test.
* Removed unused dependencies


## Dependencies

| PACKAGE                               | LICENSE |
|---------------------------------------|---------|
| **Golang**                            |         |
| github.com/labstack/echo              |   MIT   |
| github.com/jinzhu/gorm                |   MIT   |
| github.com/dgrijalva/jwt-go           |   MIT   |
| github.com/rs/zerolog                 |   MIT   |
| github.com/go-playground/validator    |   MIT   |
| github.com/nbutton23/zxcvbn-go        |   MIT   |
| github.com/stretchr/testify           |   MIT   |
| github.com/go-yaml/yaml               |  Apache |
| **node.js**                           |         |
| npmjs.com/package/multi-file-swagger  |   MIT   |

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

* run `make setup` then after that, `make serve` to use [swagger](https://github.com/johncoleman83/cerebrum#swagger) to test endpoints.

* or use `curl` to login with bootstrap data
  ```
  $ curl -X POST "http://localhost:8080/login" \
      -H "accept: application/json" \
      -H "Content-Type: application/json" \
      -d "{ \"password\": \"admin\", \"username\": \"admin\"}"
  ```

## Swagger

Visit `http://localhost:8080/swaggerui/` (with trailing slash) in your browser to test out some endpoints!

* use `make TYPE=yaml swagger` or [multi-file-swagger](https://www.npmjs.com/package/multi-file-swagger) to generate the swagger `.yaml` spec file.
* more docs on creating a swagger server: https://goswagger.io/use/spec.html#parsing-rules
* this application uses the latest distribution from [github.com/swagger-api/swagger-ui](https://github.com/swagger-api/swagger-ui/tree/master/dist)
  * this repo uses swagger docs file `full_spec.yaml` as an update from the
    upstream branch

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

## License

Currently licensed under MIT, check the [LICENSE](LICENSE.md) file for details.

## Author

[David John Coleman II](https://davidjohncoleman.com)
**[Initial Commit by GORSK](https://www.ribice.ba/refactoring-gorsk/)**
