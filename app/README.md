# Simple Go Web with Nginx
This is a Go learning project with two major directories.
## Golang app (app)

Simple Golang web application written with the net/http package. The entire app is in the `app` directory. The entire app is inside `main.go` for simplicity.

## Nginx (web)

Nginx is used as a reverse proxy when deploying the web app to a cloud provider. The `nginx.conf` alongside the `Dockerfile` is found in the `app` directory.

## Local development environment
You can use the `docker-compose.yml` file to spin up both the Golang application and Nginx web servers. But you could keep it simple by with `go run main.go` instead.
