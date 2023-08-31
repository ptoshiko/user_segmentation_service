run:
	docker-compose up --build

fclean:
	docker stop $$(docker ps -qa) \
	docker rm $$(docker ps -qa) \
	docker rmi $$(docker images -qa) \
	docker volume rm $$(docker volume ls -q) \
	docker network rm $$(docker network ls -q) \

testdb:
	docker run \
	--rm --name postgres_test \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=postgres_test \
	-p 5432:5432 \
	-d postgres:latest

t:
	go test -v