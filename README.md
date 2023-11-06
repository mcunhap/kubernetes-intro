# Introduction to Kubernetes

## Docker

To understand the basics of Kubernetes is important to know the basics of Docker. All content describe below can be found in [Docker official documentation](https://docs.docker.com/get-started/overview/)

Docker is an open platform for developing, shipping, and running applications. Docker enables you to separate your applications from your infrastructure so you can deliver software quickly. With Docker, you can manage your infrastructure in the same ways you manage your applications. By taking advantage of Docker's methodologies for shipping, testing, and deploying code, you can significantly reduce the delay between writing code and running it in production.

There is three main concepts that we must understand in Docker basics:

1. Images: an image is a read-only template with instructions for creating a Docker container. Often, an image is based on another image, with some additional customization. For example, you may build an image which is based on the ubuntu image, but installs the Apache web server and your application, as well as the configuration details needed to make your application run.  
2. Containers: a container is a runnable instance of an image. You can create, start, stop, move, or delete a container using the Docker API or CLI. You can connect a container to one or more networks, attach storage to it, or even create a new image based on its current state.  
By default, a container is relatively well isolated from other containers and its host machine. You can control how isolated a container's network, storage, or other underlying subsystems are from other containers or from the host machine.  
A container is defined by its image as well as any configuration options you provide to it when you create or start it. When a container is removed, any changes to its state that aren't stored in persistent storage disappear.  
3. Registry: a Docker registry stores Docker images. Docker Hub is a public registry that anyone can use, and Docker looks for images on Docker Hub by default. You can even run your own private registry.  

### Hands on

At this example we create a simple Golang app that runs a web server with a /hello route.
We're going to build this image and run in a docker container to understand the basics.  
ps: all files are located in `docker-example` directory

Before jump in image building we need to write a simple web server as described.  

. create a file `main.go` and copy this code into it

```go
package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Server is running on http://localhost:8080/hello")
	http.ListenAndServe(":8080", nil)
}
```

Now we can create a Dockerfile that will contains the instructions to create the image.  

. create a file `Dockerfile` and copy this code into it

```dockerfile
FROM golang:1.21.3-alpine3.18

WORKDIR /app
COPY . .
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
```

At the same directory that you have this files, you can run the `docker build` command described below to generate the `docker-example` image.

```sh
docker build -t docker-example .
```

Now you can check your images with

```sh
docker images
```

With the image we will be able to create and run a new container to execute it.

```sh
docker run --rm -p 8080:8080 docker-example
```

Note that we are exposing port 8080 from container in host 8080 (-p 8080:8080 option). Without this the web server will be unreachable from host.
The --rm flag automatically remove the container when exits.

After run the container we will be able to check basic informations about it with

```sh
docker ps
```

To access your application you can execute a simple `curl` command to `localhost:8080/hello` as described

```sh
curl localhost:8080/hello
```

We expect the following answer.

```sh
Hello World%
```

If we execute the container without exposing the port as mentioned before, than we will have some problems

```sh
docker run --rm docker-example
```

```sh
curl localhost:8080/hello
curl: (7) Failed to connect to localhost port 8080 after 5 ms: Couldn't connect to server
```

#### Push to registry

If we want to push the generated image to any registry we can execute a `docker push` command. This will allows anyone to pull this image from the registry if have access.

## Kubernetes
