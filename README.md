# GoService

This is a minimal Go HTTP service.

The `docker-compose` file includes an nginx web service for serving static content, and a backend Go HTTP app.

It is recommended to run the app using `docker-compose`. This can be achieved by running the command:

`make docker`

Which will build the containers and set nginx to listen on port 80.