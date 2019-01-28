build_dev_image:
	docker build . -f dev.dockerfile -t mindmakerdev
start_dev:
	docker run --rm -it -p 8080:8080 -v $(shell pwd):/go/src/github.com/stanleynguyen/mindmaker mindmakerdev bash