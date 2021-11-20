# Go-Proxy Sample Project

## Intro

The project tries to create a simple reverse proxy with two features. 

1. Utilize round-robin loadbalancing
2. Include an in-memory page cache.

It supports multiple clients that are able to communicate with multiple services, consisting of multiple backends each.

## Sources 
Sources used in this project:

1. [Writing a Reverse Proxy in just one line with Go](https://hackernoon.com/writing-a-reverse-proxy-in-just-one-line-with-go-c1edfa78c84b)
2. [hlts2/round-robin](https://github.com/hlts2/round-robin)
3. [Learning Go Lang Days 13 to 18](https://medium.com/codex/learning-go-lang-days-13-to-18-building-a-caching-reverse-proxy-in-go-lang-a0965495c329)
4. [ResponseWriter custom](https://stackoverflow.com/a/65895198/2766769)

The above sources used as base for my code plus education purposes to really get the idea without any assumptions 

## Prerequisites

The tools needed for the exercise are:

1. make
2. Docker/Podman

## Build
We will use make to simplify and automate the commands used in this exercise. Feel free to check them by either a) running make in your terminal or b) reading Makefile.

### Normal build process.

In order to build the project :
```
make build-app run
```

Optionally use the flags:
```
GOOS=linux GOARCH=amd64 go build -o go-proxy main.go
```

### Docker Image

A sample tag is used here `v1`
```
make build-image run-image
```

## Config

go-proxy app can be configured using a `config.yml` file providing where to listen to and the services to forward traffic to.

The structure is as follows:
```
proxy:
  listen:
    address: "127.0.0.1"
    port: 8080
  services:
    - name: mock-service-1
      domain: mock-service-1.localhost.local
      hosts:
        - address: "127.0.0.1"
          port: 8090
        - address: "127.0.0.1"
          port: 8091
        - address: "127.0.0.1"
          port: 8092
    - name: mock-service-2
      domain: mock-service-2.localhost.local
      hosts:
        - address: "127.0.0.1"
          port: 9090
        - address: "127.0.0.1"
          port: 9091
        - address: "127.0.0.1"
          port: 9092
```

## Future Improvements

* Add unit tests 
* Add integration tests
* Add proper logging with debug switch
* Perform loadtesting
* Add code versioning
* Add toggle functionality to enable/disable cache/lb_strategy
* Add health checks for the alive backends
* Add toggle functionality to purge cache and/or change cache duration
