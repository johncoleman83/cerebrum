RUNNING_CONTAINER := $(shell docker ps -aqf 'name=cerebrum_mysql_dev_db_1')

refresh: clean docker

clean:
	docker stop $(RUNNING_CONTAINER) && docker rm $(RUNNING_CONTAINER)

docker:
	docker-compose up -d
