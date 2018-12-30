# cerebrum make

refresh: clean_dev clean_test docker_dev

test: test_all clean_test

test_all:
	go test ./...

clean_test:
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_test_db_no_ | cut -f1 | xargs docker stop
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_test_db_no_ | cut -f1 | xargs docker rm

clean:
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql | cut -f1 | xargs docker stop
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql | cut -f1 | xargs docker rm

docker:
	docker-compose up -d

lint:
	golint pkg/...
	golint cmd/...