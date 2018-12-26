# cerebrum make

refresh: clean_dev clean_test docker_dev

clean_test:
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_test_no_ | cut -f1 | xargs docker stop
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_test_no_ | cut -f1 | xargs docker rm

clean_dev:
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_dev_db_1 | cut -f1 | xargs docker stop
	docker ps --all --format "{{.ID}}\t{{.Names}}" | grep cerebrum_mysql_dev_db_1 | cut -f1 | xargs docker rm

docker_dev:
	docker-compose up -d
