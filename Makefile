run:
	docker-compose up --build

fclean:
	docker stop $$(docker ps -qa) \
	docker rm $$(docker ps -qa) \
	docker rmi $$(docker images -qa) \
	docker volume rm $$(docker volume ls -q) \
	docker network rm $$(docker network ls -q) \