# GW Example
#### gRPC Gateway example project

### What is it?
___
This is a simple Echo&Reverse gRPC Server, which uses [gRPC-gateway](https://grpc-ecosystem.github.io/grpc-gateway/) to communicate also via REST.


### Usage
Make sure you have all the dependencies:

```bash
$ brew install golang
$ brew install protobuf
$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

Run it on your machine: 

```bash
$ make build
$ ./bin/service
```

If you want to use *REST*:

You just need a standard HTTP client to invoke the **echo** endpoint
    
```bash
$ curl localhost:8080/v1/example/echo -d '{"value":"Hello gRPC!"}'
{"value":"Hello gRPC"}
```
    
Or the **reverse** endpoint:

```bash
$ curl localhost:8080/v1/example/reverse -d '{"value":"Hello gRPC!"}'
   {"value":"!CPRg olleH"}
```

    

If you want to use *gRPC*:

This will build the client executable and let it communicate to the server via gRPC.
The client is cli application that currently supports the `!echo` and `!reverse` commands (use `!quit` to exit).
```bash
$ make build-client
$ ./bin/client
2019/08/14 18:13:31 Connecting to gRPC server @ 0.0.0.0:9090
Commands: !echo, !reverse, !quit
cmd> !echo
Type a message > Hello World
2019/08/14 18:13:37 client >>> Hello World
2019/08/14 18:13:37 server >>> Hello World
cmd> !reverse
Type a message > Hello World
2019/08/14 18:13:44 client >>> Hello World
2019/08/14 18:13:44 server >>> dlroW olleH
cmd> !quit
2019/08/14 18:13:48 Client shutting down...
```

### Usage with Docker

```bash
$ make docker
$ docker-compose up
```

Use a gRPC or HTTP client as showed above to interact with the server.
