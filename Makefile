# cerebrum make
DEV_CONTAINER := $(shell docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_dev_db | cut -f1)

.PHONY: list # show all make targets
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' #| sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

.PHONY: help # show all make targets with descriptions
help:
	@echo "--------------------"
	@echo "-------help---------"
	@echo "--------------------"
	@echo "TARGETS: DESCRIPTION"
	@echo "--------------------"
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1: \2/' | expand -t20

.PHONY: branch # start off a new clean git branch from origin/master head commit
branch:
ifdef NAME
	git checkout master
	git fetch origin master
	git reset --hard origin/master
	git checkout -b $(NAME)
else
	@echo 'Usage: $ make NAME=XXXXX branch'
	@echo 'where NAME is any valid new git branch'
	@echo 'run `$$ make help` for more info'
endif

.PHONY: setup # a "one click" type start from scratch to refresh the db bootstrap & serve the application
setup: refresh serve

.PHONY: refresh # refresh all docker db's and bootstrap a new dev db
refresh:
	@make ENV=all clean
	@make docker
	@make bootstrap

.PHONY: serve # starts backend server, executes $ go run cmd/api/main.go
serve:
	go run cmd/api/main.go

.PHONY: swagger # entry point to generate swagger support docs
swagger:
ifeq ($(CMD),spec)
	@make TYPE=$(TYPE) swagger_spec
else ifeq ($(CMD),client)
	@make CMD=$(CMD) TYPE=$(TYPE) swagger_ui
else ifeq ($(CMD),server)
	@make CMD=$(CMD) TYPE=$(TYPE) swagger_ui
else
	@echo 'Usage: $ make CMD=XXXXX TYPE=XXXXX swagger'
	@echo 'where CMD is the generate command and TYPE is a valid file extension like `yaml`'
	@echo 'run `$$ make help` for more info'
endif

.PHONY: swagger_spec # generate swagger spec using inline comments, use env variable TYPE to specify the output file type
swagger_spec:
ifdef TYPE
	cd cmd/api/ && \
		~/go/bin/swagger generate \
		--output=../../logs/swagger.log spec \
		--scan-models \
		--output=../../third_party/swaggerui/spec/swagger.$(TYPE)
else
	@echo 'Usage: $ make CMD=spec TYPE=XXXXX swagger'
	@echo 'where CMD is the generate command and TYPE is a valid file extension like `yaml`'
	@echo 'run `$$ make help` for more info'
endif


.PHONY: swagger_ui # generate swagger client or server side template files, must already have a valid swagger.yaml file
swagger_ui:
ifdef TYPE
	~/go/bin/swagger generate \
		--output=logs/swagger.log $(CMD) \
		--target=cmd/api/ \
		--spec=third_party/swaggerui/spec/swagger.$(TYPE) \
		--api-package=../../third_party/swaggerui/api \
		--client-package=../../third_party/swaggerui/client \
		--model-package=../../third_party/swaggerui/model \
		--server-package=../../third_party/swaggerui/server \
		--name=cerebrum
else
	@echo 'Usage: $ make CMD=XXXXX TYPE=XXXXX swagger'
	@echo 'where CMD is the generate command and TYPE is a valid file extension like `yaml`'
	@echo 'run `$$ make help` for more info'
endif

.PHONY: docker # start docker dependencies, executes $ docker-compose up
docker:
	docker-compose up --detach
	@echo "... zzz"
	@echo "going to sleep to mysql enough time to startup"
	sleep 10

.PHONY: bootstrap # bootstrap the db with dev models
bootstrap:
	go run cmd/bootstrap/main.go

.PHONY: mysql # login to mysql dev container to inspect
mysql:
	docker exec -it $(DEV_CONTAINER) mysql -u root

.PHONY: test # run all tests
test:
	@make test_go
	@make ENV=test clean
	@make lint

.PHONY: test_go # run all go file tests
test_go:
	go test ./...

.PHONY: lint # run linters on go package
lint:
	golint pkg/...
	golint cmd/...

.PHONY: clean # removes docker containers from the environmental variable ENV such as `test`
clean:
ifeq ($(ENV),dev)
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_dev_db | cut -f1 | xargs docker stop
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_dev_db | cut -f1 | xargs docker rm
else ifeq ($(ENV),test)
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_test_db_no_ | cut -f1 | xargs docker stop
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_test_db_no_ | cut -f1 | xargs docker rm
else ifeq ($(ENV),all)
	@make ENV=dev clean
	@make ENV=test clean
else
	@echo 'Usage: $ make ENV=XXXXX clean'
	@echo 'where ENV could be `test`, `dev`, `git` or `all`'
	@echo 'run `$$ make help` for more info'
endif
