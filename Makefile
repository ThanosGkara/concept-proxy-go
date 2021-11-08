

help:
	@echo "Please use 'make <target>' where <target> is one of the following:"
	@echo "  build-app 			 to build the app here"
	@echo "  build-image         to build the app container image."
	@echo "  run-image           to run the app container."

########
#  App #
########

build-app:
	GOOS=linux GOARCH=amd64 go build -o go-proxy main.go

run:
	./go-proxy

##########
# Docker #
##########
build-image:
	docker build -t go-proxy:v1 .

run-image:
	docker run --rm -it -p 8080:8080 -v ${CURDIR}/config.yml:/opt/go-proxy/config.yml --name go-proxy go-proxy:v1

.NOTPARALLEL:

.PHONY: \
	build-app\
	build-image \
	run-image \
	
