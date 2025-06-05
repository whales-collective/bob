# Docker Compose FAQ - 30 Common Questions and Answers

## 1. What is Docker Compose?

Docker Compose is a tool for defining and running multi-container Docker applications. It uses a YAML file to configure your application's services, networks, and volumes, allowing you to manage complex applications with a single command.


## 3. What is the compose.yml file?

The compose.yml file is a YAML configuration file that defines services, networks, and volumes for your multi-container application. It's the blueprint that Docker Compose uses to create and manage your application stack.

## 4. How do I start services with Docker Compose?

Use `docker composeup` to start all services defined in your compose.yml file. Add `-d` flag to run in detached mode: `docker composeup -d`

## 5. How do I stop Docker Compose services?

Use `docker composedown` to stop and remove containers, networks, and default volumes. Use `docker composestop` to just stop services without removing them.

## 6. What's the basic structure of a compose.yml file?

A basic structure includes version, services, and optionally networks and volumes:
```yaml
version: '3.8'
services:
  web:
    image: nginx
  database:
    image: postgres
```

## 7. How do I specify which Compose file to use?

Use the `-f` flag: `docker compose-f custom-compose.yml up`. You can also set the `COMPOSE_FILE` environment variable.

## 8. How do I build images with Docker Compose?

Use the `build` directive in your service definition, pointing to a directory with a Dockerfile: `build: ./my-app` or specify context and dockerfile: `build: { context: ., dockerfile: Dockerfile.dev }`

## 9. How do I set environment variables in Docker Compose?

Use the `environment` key in your service, or reference an `.env` file with `env_file`. Example: `environment: - NODE_ENV=production` or `env_file: - .env`

## 10. How do I map ports in Docker Compose?

Use the `ports` directive: `ports: - "8080:80"` maps host port 8080 to container port 80. You can also use `- "80"` for random host port assignment.

## 11. How do I define volumes in Docker Compose?

Use the `volumes` directive under services: `volumes: - ./data:/app/data` for bind mounts, or `- db_data:/var/lib/mysql` for named volumes (which must be declared at the top level).

## 12. How do I scale services with Docker Compose?

Use `docker compose up --scale service_name=3` to run 3 instances of a service, or define `deploy.replicas: 3` in your compose file (for swarm mode).

## 13. How do I view logs from Docker Compose services?

Use `docker composelogs` to see all service logs, or `docker composelogs service_name` for a specific service. Add `-f` to follow logs in real-time.

## 14. What are Docker Compose networks?

Networks allow services to communicate with each other. Docker Compose automatically creates a default network for your application, and services can reach each other using service names as hostnames.

## 15. How do I create custom networks in Docker Compose?

Define networks at the top level and reference them in services:
```yaml
networks:
  frontend:
  backend:
services:
  web:
    networks: [frontend]
```

## 16. How do I depend on other services in Docker Compose?

Use `depends_on` to specify service dependencies: `depends_on: - database`. This ensures the database starts before the web service, though it doesn't wait for the service to be ready.

## 17. How do I restart services with Docker Compose?

Use `docker composerestart` to restart all services, or `docker composerestart service_name` for a specific service. You can also set restart policies with `restart: always`.

## 18. How do I run one-off commands with Docker Compose?

Use `docker composerun service_name command` to run a command in a new container, or `docker composeexec service_name command` to run in an existing container.

## 19. What are profiles in Docker Compose?

Profiles allow you to selectively start services. Define profiles on services and use `--profile` flag: `docker compose--profile debug up` to start only services with the debug profile.

## 20. How do I override Docker Compose configurations?

Use multiple compose files: `docker compose-f compose.yml -f docker-compose.override.yml up`. The override file extends or overrides the base configuration.

## 21. How do I set health checks in Docker Compose?

Use the `healthcheck` directive:
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 3
```

## 22. How do I use secrets in Docker Compose?

Define secrets at the top level and reference them in services:
```yaml
secrets:
  db_password:
    file: ./db_password.txt
services:
  app:
    secrets: [db_password]
```

## 23. How do I limit resources for services?

Use deploy configuration (mainly for swarm mode):
```yaml
deploy:
  resources:
    limits:
      cpus: '0.5'
      memory: 512M
```

## 24. What's the difference between docker composeup and docker composestart?

`docker composeup` creates and starts containers, while `docker composestart` only starts existing stopped containers. Use `up` for initial deployment and `start` to resume stopped services.

## 25. How do I remove volumes when stopping Docker Compose?

Use `docker composedown -v` to remove named volumes along with containers and networks. Be careful as this will delete data stored in volumes.

## 26. How do I use external networks in Docker Compose?

Reference an existing network with `external: true`:
```yaml
networks:
  existing_network:
    external: true
services:
  app:
    networks: [existing_network]
```

## 27. How do I extend services in Docker Compose?

Use the `extends` keyword to inherit configuration from another service in the same or different file:
```yaml
web:
  extends:
    file: common.yml
    service: webapp
```

## 28. How do I use Docker Compose with different environments?

Create separate compose files for each environment (docker-compose.prod.yml, docker-compose.dev.yml) and use them with the `-f` flag, or use environment variables and .env files.

## 29. How do I debug Docker Compose issues?

Use `docker composeconfig` to validate and view the final configuration, `docker composeps` to check service status, and `docker composelogs` to examine service output. Add `--verbose` for detailed information.

## 30. What are some Docker Compose best practices?

Use specific image tags instead of 'latest', organize services logically, use .env files for environment-specific values, implement health checks, use multi-stage builds for smaller images, and separate development and production configurations.