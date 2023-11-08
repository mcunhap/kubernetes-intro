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

### Hands On

At this example we create a simple Golang app that runs a web server with a /hello route.
We're going to build this image and run in a docker container to understand the basics.  
ps: all files are located in `docker-example` directory

Ensure that you have Docker installed, if not you can download [here](https://www.docker.com/products/docker-desktop/).  

#### Create simple web server

Before jump into image building we need to write a simple web server as described.  

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

##### Create Dockerfile

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

At the same directory that you have this files, you can run the `docker build` command described below to generate the `go-web-app` image.

```sh
docker build -t go-web-app .
```

You can check your images with

```sh
docker images
```

#### Run a container

With the image we will be able to create and run a new container to execute it.

```sh
docker run --rm --name my-go-web-app -p 8080:8080 go-web-app
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
docker run --rm --name my-go-web-app go-web-app
```

```sh
curl localhost:8080/hello
curl: (7) Failed to connect to localhost port 8080 after 5 ms: Couldn't connect to server
```

#### Running a shell inside an existing container

It is possible to access the shell inside the container that is running your application with the following command:

```sh
docker exec -it my-go-web-app /bin/sh
```

The process that is running the shell have the same Linux namespaces as the main container proccess. This make possible to us explore the container and see how our app is running.  
Each container has an isolated process tree, and also an isolated filesystem. Checking container processes will ensure to us the process tree isolation.

```sh
/app # ps -ef
PID   USER     TIME  COMMAND
    1 root      0:00 ./main
   11 root      0:00 /bin/sh
   17 root      0:00 ps -ef
```

Listing contents from root directory inside the container show to us about filesystem isolation.  

```sh
/ # ls
app    dev    go     lib    mnt    proc   run    srv    tmp    var
bin    etc    home   media  opt    root   sbin   sys    usr
```

#### Stop and delete containers

When we need to stop a container we can execute

```sh
docker stop my-go-web-app
```

If you run the container without `--rm` option you can delete the container with

```sh
docker rm my-go-web-app
```

To make our image accessible to any user and application that we want to, then we need to push the generated image to any registry. To make easier we will push to [Docker Hub](https://hub.docker.com/).  
Before push we need to create an account in Docker Hub and tag our image following Docker Hub's rules, that is include you Docker Hub ID in the tag like below

```sh
docker tag go-web-app mcunhap/go-web-app
```

You can change `mcunhap` with your personal Docker Hub ID.  

Before push you need to login to docker hub with `docker login`, then you can run the command below to push the image

```sh
docker push mcunhap/go-web-app
```

## Kubernetes

Kubernetes is an open-source platform that enables the automation, management, and scalability of containerized applications, commonly used with Docker.  

The architecture of Kubernetes, as described in [documentation](https://kubernetes.io/pt-br/docs/concepts/overview/components/) can be seen in figure below. It consists of different components that work together to provide a resilient and scalable environment. When using Kubernetes, we deal with a Kubernetes cluster, which is a set of processing servers called nodes. Pods will be hosted on these servers, where the application will run. There are two types of nodes: the Master Node, responsible for making scalability decisions and managing cluster resources, and the Worker Nodes, responsible for hosting applications. Every Kubernetes cluster has at least one Worker Node.  

![image](docs/images/k8s-architecture.png)

The Kubernetes Master Node consists of three main components:

- **API Server:** This component exposes the Kubernetes API, serving as the entry point for all services.

- **Scheduler:** Responsible for monitoring newly created pods that have not been assigned to a node and selecting a suitable node for their execution.

- **Controller Manager:** Tasked with managing the cluster's controllers, ensuring proper cluster state maintenance.

The Kubernetes Worker Node also comprises three main components:

- **Kubelet:** An agent running on each node, ensuring containers are running in a pod.

- **Kube Proxy:** A network proxy responsible for maintaining network rules and enabling communication with pods both inside and outside the cluster.

- **Container Runtime:** The software responsible for running containers, such as Docker.


Finally, the Kubernetes architecture includes Etcd, a consistent and highly available key-value store used for all cluster data.  

### Concepts

Next, let's define the concept of a *pod*, understand what it represents, and explore the available options for *workload resources* and services in Kubernetes.

In the context of Kubernetes, the [pod](https://kubernetes.io/docs/concepts/workloads/pods/) represents the smallest deployable unit under our control for creation and management. In its typical configuration, a *pod* consists of a single container; however, in situations involving applications composed of multiple interconnected containers, the *pod* can encompass several containers. This possibility arises due to the sharing of network and storage resources among all containers contained in the *pod*. It's worth noting that *pods* in Kubernetes are ephemeral by nature, meaning they can be created, deleted, and recreated dynamically to adapt to application demands.

Normally, *pods* are not created directly as a resource in Kubernetes but are instead created from *workload resources*. A [workload resource](https://kubernetes.io/docs/concepts/workloads/) is basically an automated way of managing a set of *pods*. Kubernetes provides us with some *workload resources*:

- **Deployment**: Manages stateless applications where the state is not saved, and each transaction can be interpreted as starting from scratch.
- **StatefulSet**: Manages stateful applications where state persistence is required.
- **DaemonSet**: Ensures a specific copy of a *pod* runs on each node of the cluster. In case new nodes are added to the cluster, this *workload* guarantees that a new *pod* will appear on the new node of the cluster.
- **Job and CronJob**: *Job* provides a way to execute tasks that run once and are completed. For scenarios where the same *Job* needs to be executed multiple times according to a schedule, we can use *CronJob*.

To expose an application from a set of *pods* as a network service, Kubernetes provides the *Service* resource. With it, we can define a logical set of *pods* and a stable way to access them. The *service* makes it possible to expose this set not only internally to the cluster but also externally if necessary. Kubernetes provides us with the following types of *service*:

- **ClusterIP**: Allocates an internal IP, making the set of *pods* accessible only internally within the cluster.
- **NodePort**: Exposes the service on a fixed port of each cluster node, making it accessible from the cluster node's IP on the configured port.
- **LoadBalancer**: Exposes the service externally using an external load balancer. Note that Kubernetes does not provide a load balancer component; you need to provide one, commonly offered in public cloud clusters such as GCP.
- **ExternalName**: Maps the service to a specified Domain Name System (DNS).

### Autoscaling

Autoscaling aims to dynamically adjust processing resources based on the application's demand. This approach ensures that more resources are allocated during peak demand and deallocated during idle times. Thus, it guarantees high availability for the application and optimization of the resources used.

The autoscaling process consists of four distinct phases. The first phase is monitoring, where a system is used to collect data about the monitored resources. Next, this data is analyzed by the autoscaling system, which uses a predefined rule to plan an action to be executed. The granularity and quality of the monitoring data have a direct impact on the performance of the autoscaling system, so it is crucial that the collected data for analysis is reliable and readily available.

Kubernetes provides three possibilities for autoscaling:

- **Horizontal Pod Autoscaler (HPA):** Adds new pods that run the same application, balancing the load across more processing units. Since only new pods need to be added, there is no need to restart the running application, making HPA attractive for applications requiring high availability.

- **Vertical Pod Autoscaler (VPA):** Directly changes the resources of existing pods. As it is necessary to modify the available resources in existing pods, a service restart is required, hindering the use of VPA for applications requiring high availability.

- **Cluster Autoscaler (CA):** Increases the number of nodes in the Kubernetes cluster when it is no longer possible to allocate pods on existing nodes. It is commonly used by public clouds that provide Kubernetes clusters, such as GCP (Google Cloud Platform).

### Hands On

At this example we will run a Kubernetes single node cluster using [Minikube](https://minikube.sigs.k8s.io/docs/), that is a tool to quickly sets up a local Kubernetes cluster and is great for test and develop apps locally. Then, we will deploy the simple web app that we build in Docker section.
To install minikube you can follow instructions in minikube [documentation](https://minikube.sigs.k8s.io/docs/start/).

#### Cluster setup

Once you have installed minikube you can easily start up the Kubernetes cluster with

```sh
minikube start
```

We will also need kubectl CLI to comunicate with Kubernetes cluster and run commands. If needed you can follow this (documentation)[https://kubernetes.io/docs/tasks/tools/] to install it.

When the cluster finish start up process you can check if your cluster is working with

```sh
kubectl cluster-info
```

We expect an output like the one below

```sh
Kubernetes control plane is running at https://127.0.0.1:59579
CoreDNS is running at https://127.0.0.1:59579/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

If you check docker containers with `docker ps` you will realize that we have a new container that is running minikube. Inside this container is running our Kubernetes cluster, that is totally decoupled from our OS.

#### Deploy web application

Now we will deploy the image that we build in Docker section to Kubernetes in a way to make it accessible outside the cluster. To create the web server, build the image and push to Docker Hub you can check Docker section, if haven't done yet.

The easiest way to deploy our application is using `kubectl create` command like this:

```sh
kubectl create deployment go-web-app --image=mcunhap/go-web-app --port=8080
```

Another option is to create an `yaml` file and describe exactly what we want. To do it in that way you can create a file named `go-web-app.yml` and add the following code to it

```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-web-app
  labels:
    app: go-web-app
spec:
  selector:
    matchLabels:
      app: go-web-app
  replicas: 1
  template:
    metadata:
      labels:
        app: go-web-app
    spec:
      containers:
      - name: go-web-app
        image: mcunhap/go-web-app
        resources:
          requests:
            cpu: "1000m"
            memory: "1Gi"
          limits:
            cpu: "1000m"
            memory: "1.5Gi"
        ports:
          - containerPort: 8080
            protocol: TCP
```

Note that in this yaml file we are specifying more information, like the resources that each container can use and the number of desired replicas.  
With this yaml written we can create the deployment running

```sh
kubectl apply -f go-web-app.yml
```

We can check our deployments with

```sh
kubectl get deployments
```

And our pods with

```sh
kubectl get pods
```

To get more informations about pods you can add `-o wide` to `kubectl get pods` command, like

```sh
kubectl get pods -o wide
```

Now we have our image deployed to the Kubernetes, but how can we access it? Each pod gets its own IP, but this address is internal to the cluster and isn't accessible from outside of it. So, if we try to make

```
curl localhost:8080/hello
```

will result in a fail, since the application is not accessible outside the cluster. To make accessible we need to expose it through a Service.  
Now we will create an Service object of type NodePort. Create a file named `go-web-app-service.yml` and add the following code to it.

```yml
apiVersion: v1
kind: Service
metadata:
  name: go-web-app-service
spec:
  type: NodePort
  selector:
    app: go-web-app
  ports:
  - name: gateway
    port: 8080
    targetPort: 8080
    protocol: TCP
```

Here we are defining a NodePort service and mapping the actual port on which our application is running inside the container (targetPort) to Service specific port (port) expliciting TCP protocol. Note that we are mapping this service to
app `go-web-app`, so Kubernetes know to which app this Service is pointing to. Then, we can apply the service with

```sh
kubectl apply -f go-web-app-service.yml
```

To check our services we can run

```sh
kubectl get services
```

If we check the output we will see that Service port 8080 is mapped to 32621 in Kubernetes node. Note that this port is selected randomly by default, so it might change for each one.

```sh
NAME                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
go-web-app-service   NodePort    10.100.240.92   <none>        8080:32621/TCP   9m21s
```

Now, if we try to access our application we still get a failure...

```sh
curl localhost:32621/hello
curl: (7) Failed to connect to localhost port 32621 after 7 ms: Couldn't connect to server
```

But, why? Because minikube is not running in our host, but inside a docker container. So, we have two options:

1. Access the container that is running minikube and access our application in port 32621
2. Access our application from our host with minikube tunnel


To make the first option we can run

```sh
docker exec -it minikube /bin/sh
```

Inside the container just make the request to `localhost:32621/hello` and...

```sh
curl localhost:32621/hello
Hello World
```

To access using minikube tunnel, you can execute

```sh
minikube service go-web-app-service
```

Then will be mapped another port to use locally and you can access in your browser, or make the request with curl

```sh
curl localhost:<port-selected>/hello
```

#### Scale our application horizontally

Now that we have our application running and accessible we can manage to increase and decrease the replicas. First we will make it manually by editing our `go-web-app.yml`
file changing the replicas field and applying changes with `kubectl apply -f go-web-app.yml` command:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-web-app
  labels:
    app: go-web-app
spec:
  selector:
    matchLabels:
      app: go-web-app
  replicas: 3
  template:
    metadata:
      labels:
        app: go-web-app
    spec:
      containers:
      - name: go-web-app
        image: mcunhap/go-web-app
        resources:
          requests:
            cpu: "1000m"
            memory: "1Gi"
          limits:
            cpu: "1000m"
            memory: "1.5Gi"
        ports:
          - containerPort: 8080
            protocol: TCP
```

Or just running the command

```sh
kubectl scale deployment go-web-app --replicas=3
```

After apply the upscale we can check our pods with `kubectl get pods` and we will see something like:

```sh
NAME                         READY   STATUS    RESTARTS   AGE
go-web-app-59d87f9d7-dzqlt   1/1     Running   0          40m
go-web-app-59d87f9d7-kbb7s   1/1     Running   0          83s
go-web-app-59d87f9d7-p8gqw   1/1     Running   0          83s
```

If we want to scale down we can follow the same process, but setting the replicas number to a lower value than current value.

#### HPA Example

Now we want to scale our application automatically. To see an example we are going to follow Kubernetes [documentation](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/) HPA walkthrough.

First we need to enable metrics-server in our minikube cluster. To do so, just run

```sh
minikube addons enable metrics-server
```

[Metrics server](https://github.com/kubernetes-sigs/metrics-server) is a source of container resource metrics for Kubernetes built-in autoscale pipelines. Metrics Server collects resource metrics from Kubelets and exposes them in Kubernetes apiserver through Metrics API for use by Horizontal Pod Autoscaler and Vertical Pod Autoscaler.

Now that we have enabled metrics-server we will be able to get resource metrics from pods and make decisions based on them. First we will create a file name `php-apache.yml` and include the following code to it

```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: php-apache
spec:
  selector:
    matchLabels:
      run: php-apache
  template:
    metadata:
      labels:
        run: php-apache
    spec:
      containers:
      - name: php-apache
        image: registry.k8s.io/hpa-example
        ports:
        - containerPort: 80
        resources:
          limits:
            cpu: 500m
          requests:
            cpu: 200m
---
apiVersion: v1
kind: Service
metadata:
  name: php-apache
  labels:
    run: php-apache
spec:
  ports:
  - port: 80
  selector:
    run: php-apache
```

Here we are describing a deployment and a service in the same yaml file (yes! that is possible!), this deployment will use an example image from kubernetes registry and each container will request 200m CPU with a limit of 500m. That means that in the beginning only 200m CPU will be requested, but if necessary we can go up to 500m, basically. We're going to expose this deploy with a ClusterIP service, that means that this will only be accessible inside the cluster.

Now we can apply this to our cluster with

```sh
kubectl apply -f php-apache.yml
```

We can check our pods and service with the commands that we already learn.

Now, we need to set up our HPA. So, we are going to create a file named `php-apache-hpa.yml` and include the following code to it

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: php-apache-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: php-apache
  minReplicas: 1
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 50
```

Here we are creating a HPA with minimum replicas 1 and maximum replicas 10, and the rule to scale will be based on CPU utilization with an average of 50. That means that our applicatio will scale up when the average utilization of pods CPU reach 50. Note that we are binding this HPA to php-apache deployment.

To apply this HPA we can run

```sh
kubectl apply -f php-apache-hpa.yml
```

We can check our HPA status with

```sh
kubectl get hpa php-apache-hpa
```

With everything working we can increase load of our application. To do this we will start a different pod to act as a client. The container that will act like a client will contain a infinite loop sending queries to the php-apache server. Run the code below in a separate terminal

```sh
kubectl run -i --tty load-generator --rm --image=busybox:1.28 --restart=Never -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://php-apache; done"
```

To continually watch HPA status we can run in another terminal

```sh
kubectl get hpa php-apache-hpa -w
```

In a few minutes we will see the higher CPU load, and then more replicas getting added.

To stop sending load to the application just hit `<Ctrl> + C` to the terminal that is running our load container. In a few minutes we can check the load decreasing and then php-apache replicas decreasing also.

### ref
[https://github.com/knrt10/kubernetes-basicLearning](https://github.com/knrt10/kubernetes-basicLearning)
