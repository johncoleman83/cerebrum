# cerebrum make
DEV_CONTAINER := $(shell docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_dev_db | cut -f1)
TEST_CONTAINER := $(shell docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_test_db | cut -f1)

.PHONY: list # show all make targets
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' #| sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

.PHONY: help # show all make targets with descriptions
help:
	@echo "-----------------------"
	@echo "|         help        |"
	@echo "-----------------------"
	@echo "| TARGET: DESCRIPTION |"
	@echo "-----------------------"
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1: \2/'

.PHONY: godoc # run godoc server and site on port 6060 to see package docs
godoc:
	godoc -play=true -index -http=:6060

.PHONY: init # install package and development dependencies such as docker, npm, swagger compiler, golang/dep
init:
	./scripts/init.py

.PHONY: dep # install the project's dependencies
dep:
	dep ensure -v

.PHONY: update # update all go package dependencies
update:
	dep ensure --update -v

.PHONY: setup # start docker dev db, bootstrap it and serve application
setup:
	@make ENV=dev docker
	@make bootstrap
	go run cmd/api/main.go

.PHONY: serve # starts backend server, executes $ go run cmd/api/main.go
serve:
	@make ENV=dev docker
	go run cmd/api/main.go

.PHONY: docker # start docker dev dependencies, executes $ docker-compose --file ./configs/docker/docker-compose.yml --file ./configs/docker/docker-compose.dev.yml up --detach
docker:
ifeq ($(ENV),dev)
	if ! [ -z $(TEST_CONTAINER) ]; then	\
		if [[ `docker inspect -f {{.State.Running}} $(TEST_CONTAINER)` = true ]]; then docker stop $(TEST_CONTAINER); fi \
	fi
	if [ -z $(DEV_CONTAINER) ] || [ `docker inspect -f {{.State.Running}} $(DEV_CONTAINER)` = false ]; then	\
		docker-compose --file ./configs/docker/docker-compose.yml up --detach db_dev \
		&& echo -e "... zzz\ngoing to sleep to allow mysql enough time to startup" \
		&& sleep 10; fi
else ifeq ($(ENV),test)
	if ! [ -z $(DEV_CONTAINER) ]; then	\
		if [[ `docker inspect -f {{.State.Running}} $(DEV_CONTAINER)` = true ]]; then docker stop $(DEV_CONTAINER); fi \
	fi
	if [ -z $(TEST_CONTAINER) ] || [ `docker inspect -f {{.State.Running}} $(TEST_CONTAINER)` = false ]; then	\
		docker-compose --file ./configs/docker/docker-compose.yml up --detach db_test \
		&& echo -e "... zzz\ngoing to sleep to allow mysql enough time to startup" \
		&& sleep 10; fi
else
	@echo 'Usage: $ make ENV=XXXXX docker'
	@echo 'where ENV could be `test` or `dev`'
	@echo 'run `$$ make help` for more info'
endif

.PHONY: bootstrap # bootstrap the db with dev models
bootstrap:
	go run scripts/bootstrap/main.go

.PHONY: mysql # login to mysql dev container to inspect
mysql:
	@make ENV=dev docker
	docker exec -it $(DEV_CONTAINER) mysql -u root

.PHONY: test_script # runs a test script $ go run scripts/testing/main.go
test_script:
	@make ENV=dev docker
	go run scripts/testing/main.go

.PHONY: test # run all tests
test:
	@make ENV=test docker
	@make test_go
	@make lint

.PHONY: test_go # run all go file tests
test_go:
	go test `go list ./... | grep -v -e pkg/utl/mock`

.PHONY: lint # run linters on go package
lint:
	golint pkg/...
	golint cmd/...

.PHONY: swagger # compile swagger spec file from swagger directories. Usage: make TYPE=XXXXX swagger
swagger:
ifdef TYPE
	cd third_party/swaggerui/spec && \
		multi-file-swagger -o $(TYPE) index.yaml > compiled/full_spec.$(TYPE) && \
		cp -rf compiled/full_spec.$(TYPE) ../dist/full_spec.yaml && \
		cd -
else
	@echo 'Usage: $ make TYPE=XXXXX swagger'
	@echo 'where TYPE is a valid file extension like `yaml` or `json`'
	@echo 'run `$$ make help` for more info'
endif

.PHONY: clean # removes docker containers from the environmental variable ENV such as `test`
clean:
ifeq ($(ENV),dev)
	docker stop $(DEV_CONTAINER)
	docker rm $(DEV_CONTAINER)
else ifeq ($(ENV),test)
	docker stop $(TEST_CONTAINER)
	docker rm $(TEST_CONTAINER)
else ifeq ($(ENV),all)
	@make ENV=dev clean
	@make ENV=test clean
else
	@echo 'Usage: $ make ENV=XXXXX clean'
	@echo 'where ENV could be `test`, `dev`, `git` or `all`'
	@echo 'run `$$ make help` for more info'
endif

.PHONY: refresh # refresh all docker db's and bootstrap a new dev db
refresh:
	@make ENV=all clean
	@make ENV=dev docker
	@make bootstrap

.PHONY: remove_mysql_image # removes mysql:latest docker image
remove_mysql_image:
	docker images --format "{{.Repository}}:{{.Tag}}\t{{.ID}}" --all | grep mysql:latest | cut -f2 | xargs docker rmi
