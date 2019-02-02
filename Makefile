build_dev_image:
	docker build . -f dev.dockerfile -t mindmakerdev
start_dev:
	mkdir -p tmp/redis/data && \
	docker network create mindmaker-net || true && \
	docker stop mindmakerredis || true && \
	docker rm mindmakerredis || true && \
	docker run -td --network mindmaker-net --name mindmakerredis -v $(shell pwd)/tmp/redis/data:/data redis && \
	docker run --rm -it --network mindmaker-net -p 8080:8080 -v $(shell pwd):/go/src/github.com/stanleynguyen/mindmaker mindmakerdev bash
