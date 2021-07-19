# Go REST API Course through Containers and Kubernetes.

This is code from Elliot Forbes excellent Go REST API Course found on:

https://tutorialedge.net/ and at: https://github.com/TutorialEdge/go-rest-api-course

The files here were typed by hand while following the videos and in some cases modified.

They were tested during the development process but may contain bugs.

There are also som additional scripts for running and testing the examples.

Some of the scripts in this directory expect a shell file ENV.sh
that contains the following environment variables.

Tailor these values to your environment.

```
export DB_USERNAME=postgres
export DB_PASSWORD=postgres
export DB_HOST=192.168.0.100
export DB_TABLE=postgres
export DB_PORT=5432
```

To run the K8S example with a local docker hosted database see: docs/RunWithLocalDB.md
