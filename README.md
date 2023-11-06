# Introduction to kubernetes with a simple example

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

#### Build Image

```sh
docker build -t docker-example .
```

#### Execute image

```sh
docker run --rm -p 8080:8080 docker-example
```

Note that we are exposing port 8080 from container in host 8080. Without this the web server will be unreachable from host.
The --rm flag automatically remove the container when exits.

#### Check all container that is running

```sh
docker ps
```

#### Run a request to docker-example container

```sh
curl localhost:8080/hello
```

Running correctly we expect the following response

```sh
Hello World%
```

If we execute the container without exposing the port as mentioned before, than we expect the following response

```sh
docker run --rm docker-example

curl localhost:8080/hello
curl: (7) Failed to connect to localhost port 8080 after 5 ms: Couldn't connect to server
```

#### Push to registry

If we want to push the generated image to any registry we can execute a `docker push` command. This will allows anyone to pull this image from the registry if have access.

## Kubernetes
