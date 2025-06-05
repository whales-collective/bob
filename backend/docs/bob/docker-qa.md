# Docker FAQ - 30 Common Questions and Answers

## 1. What is Docker?

Docker is a containerization platform that allows you to package applications and their dependencies into lightweight, portable containers. These containers can run consistently across different environments, from development to production.

## 2. What's the difference between Docker containers and virtual machines?

Containers share the host OS kernel and are more lightweight, while VMs include a full guest operating system. Containers start faster, use fewer resources, and provide better performance than VMs for most applications.

## 3. How do I install Docker?

Download Docker Desktop from docker.com for Windows/Mac, or install Docker Engine on Linux using your package manager. For Ubuntu: `sudo apt-get update && sudo apt-get install docker.io`

## 4. What is a Docker image?

A Docker image is a read-only template containing the application code, runtime, libraries, and dependencies needed to run an application. Images are used to create containers.

## 5. What is a Docker container?

A Docker container is a running instance of a Docker image. It's an isolated, executable package that includes everything needed to run an application.

## 6. How do I pull an image from Docker Hub?

Use the `docker pull` command followed by the image name: `docker pull nginx` or `docker pull ubuntu:20.04` for a specific version.

## 7. How do I run a container?

Use `docker run` followed by the image name: `docker run nginx` or `docker run -d -p 8080:80 nginx` to run in detached mode with port mapping.

## 8. How do I list running containers?

Use `docker ps` to see running containers, or `docker ps -a` to see all containers (including stopped ones).

## 9. How do I stop a container?

Use `docker stop <container_id>` or `docker stop <container_name>`. You can also use `docker kill` for forceful termination.

## 10. How do I remove a container?

Use `docker rm <container_id>` to remove a stopped container, or `docker rm -f <container_id>` to force remove a running container.

## 11. What is a Dockerfile?

A Dockerfile is a text file containing instructions to build a Docker image. It defines the base image, dependencies, commands, and configuration needed for your application.

## 12. How do I build an image from a Dockerfile?

Use `docker build -t <image_name> .` in the directory containing your Dockerfile. The `-t` flag tags the image with a name.

## 13. How do I map ports between host and container?

Use the `-p` flag with `docker run`: `docker run -p 8080:80 nginx` maps host port 8080 to container port 80.

## 14. How do I mount volumes in Docker?

Use the `-v` flag: `docker run -v /host/path:/container/path nginx` or use named volumes: `docker run -v myvolume:/container/path nginx`

## 15. What's the difference between COPY and ADD in Dockerfile?

`COPY` simply copies files from host to container. `ADD` has additional features like extracting tar files and downloading URLs, but `COPY` is preferred for simple file copying.

## 16. How do I access a running container's shell?

Use `docker exec -it <container_id> /bin/bash` or `docker exec -it <container_id> sh` if bash isn't available.

## 17. How do I view container logs?

Use `docker logs <container_id>` to view logs, or `docker logs -f <container_id>` to follow logs in real-time.


## 20. How do I list all Docker images?

Use `docker images` or `docker image ls` to see all locally stored images with their tags, sizes, and creation dates.

## 21. How do I remove Docker images?

Use `docker rmi <image_id>` to remove an image, or `docker image prune` to remove unused images. Use `docker image prune -a` to remove all unused images.

## 22. What are Docker networks?

Docker networks allow containers to communicate with each other and external systems. Docker provides bridge, host, overlay, and custom network types.

## 23. How do I create a custom Docker network?

Use `docker network create <network_name>` to create a bridge network, or specify the driver: `docker network create -d bridge mynetwork`

## 24. How do I connect containers to a network?

Use `--network` flag when running: `docker run --network mynetwork nginx` or connect existing containers: `docker network connect mynetwork <container_id>`

## 25. What is the difference between CMD and ENTRYPOINT?

`CMD` provides default arguments that can be overridden, while `ENTRYPOINT` defines the main command that always runs. They can be used together where `ENTRYPOINT` sets the command and `CMD` provides default arguments.

## 26. How do I set environment variables in containers?

Use the `-e` flag: `docker run -e MY_VAR=value nginx` or use `--env-file` to load from a file: `docker run --env-file .env nginx`

## 27. How do I inspect a container or image?

Use `docker inspect <container_id>` for containers or `docker inspect <image_name>` for images to get detailed JSON information about configuration and metadata.

## 28. What is Docker Hub?

Docker Hub is a cloud-based registry service where you can store, share, and manage Docker images. It's the default registry for pulling public images.

## 29. How do I push an image to Docker Hub?

First tag your image: `docker tag myimage username/myimage:tag`, then push: `docker push username/myimage:tag`. You need to be logged in with `docker login`.

## 30. How do I clean up Docker resources?

Use `docker system prune` to remove unused containers, networks, and images. Add `-a` flag to also remove unused images. Use `docker volume prune` for unused volumes.