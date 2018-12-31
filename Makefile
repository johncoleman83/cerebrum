# cerebrum make
DEV_CONTAINER := $(shell docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_dev_db | cut -f1)

.PHONY: list # show all make targets
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' #| sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

.PHONY: help # show all make targets with descriptions
help:
	@echo "TARGETS: DESCRIPTION"
	@echo "--------------------"
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1: \2/' | expand -t20

.PHONY: setup # a "one click" type start from scratch to refresh the db bootstrap & serve the application
setup: refresh serve

.PHONY: refresh # refresh all docker db's and bootstrap a new dev db
refresh: clean_dev clean_test docker bootstrap

.PHONY: serve # starts backend server, executes $ go run cmd/api/main.go
serve:
	go run cmd/api/main.go

.PHONY: docker # start docker dependencies, executes $ docker-compose up
docker:
	docker-compose up --detach
	@echo "... zzz"
	@echo "going to sleep to mysql enough time to startup"
	sleep 10

.PHONY: bootstrap # bootstrap the db with dev models
bootstrap:
	go run cmd/bootstrap/main.go

.PHONY: bootstrap # login to mysql dev container to inspect
mysql:
	docker exec -it $(DEV_CONTAINER) mysql -u root

.PHONY: test # run all tests
test: test_go clean_test lint

.PHONY: test_go # run all go file tests
test_go:
	go test ./...

.PHONY: lint # run linters on go package
lint:
	golint pkg/...
	golint cmd/...

.PHONY: clean_all # clean all test and dev docker containers
clean_all: clean_test clean_dev

.PHONY: clean_dev # remove all dev docker containers
clean_dev:
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_dev_db | cut -f1 | xargs docker stop
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_dev_db | cut -f1 | xargs docker rm

.PHONY: clean_test # remove all test docker containers
clean_test:
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_test_db_no_ | cut -f1 | xargs docker stop
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_test_db_no_ | cut -f1 | xargs docker rm
