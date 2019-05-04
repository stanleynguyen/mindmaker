build_heroku_deployment:
	docker build . -f heroku.dockerfile -t registry.heroku.com/mindmaker/web
start_dev:
	docker-compose -f docker-compose.dev.yml up -d
stop_dev:
	docker-compose -f docker-compose.dev.yml down
go_in_dev_container:
	docker exec -it mindmaker_backend_1 bash
rebuild_dev:
	docker-compose -f docker-compose.dev.yml build
