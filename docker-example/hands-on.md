## Hands on

Build Image

```sh
docker build -t docker-example .
```

Execute image

```sh
docker run --rm -p 8080:8080 docker-example
```

Note that we are exposing port 8080 from container in host 8080. Without this the web server will be unreachable from host.
The --rm flag automatically remove the container when exits.

Check all container that is running

```sh
docker ps
```

Run a request to docker-example container

```sh
curl localhost:8080/hello
```
