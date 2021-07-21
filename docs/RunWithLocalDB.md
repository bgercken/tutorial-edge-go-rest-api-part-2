# How to run the K8s example with a local Docker Database

The K8s example for the course uses a hosted database for testing.

This README describes how to use the local Docker database used in the previous examples.

This example was tested on a Linux server running a local cluster. It has not been tested anywhere else. 

The steps fit into the course around the CONTAINERS AND KUBERNETES section after completing the "Defining our K8s Deployment" video.

---

### TLDR - 

The quick description is that you can use the existing docker-compose approach to start the database (from the earlier examples) and then modify the deployment and the service to use a different port (so that it does not conflict with the docker-compose app container that is running).

### Deep Dive

There are 3 basic steps to this process.

1. Change the `containerPort` in the `config/deploment.yml` file.

    For this example we change the containerPort from `8080` to `8888`.

    ```
    ---
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: comments-api
    spec:
      replicas: 1
      strategy:
        type: RollingUpdate
        rollingUpdate:
          maxSurge: 1
          maxUnavailable: 0
      selector:
        matchLabels:
          name: comments-api
      template:
        metadata:
          labels:
            name: comments-api
        spec:
          containers:
          - name: application
            image: "bgercken/comments-api:latest"
            imagePullPolicy: Always
            ports:
              - containerPort: 8888
            env:
              - name: DB_PORT
                value: "$DB_PORT"
              - name: DB_HOST
                value: "$DB_HOST"
              - name: DB_PASSWORD
                value: "$DB_PASSWORD"
              - name: DB_TABLE
                value: "$DB_TABLE"
              - name: DB_USERNAME
                value: "$DB_USERNAME"
              - name: SSL_MODE
                value: "disable"
    ```

2. Change the port in the `config/service.yml` file. 

    For this example we change the port from `8080` to `8888`.

    ```
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: comments-api2
    spec:
      type: NodePort
      selector:
        name: comments-api
      ports:
      - protocol: TCP
        port: 8888
        targetPort: 8080
    ```

3. Change the syntax of the `port-forward` command to point to the new port (8888).

    ```
    kubectl port-forward --address 0.0.0.0 service/comments-api 8888:8080
    Forwarding from 0.0.0.0:8888 -> 8080
    Handling connection for 8888
    ```

    Here we are using an additionl argument `--address 0.0.0.0` to allow the connection to bind on all interfaces. This is to faciliate testing from external hosts and is not necessary if you are running everything locally. Note that you may need to open port 8888/tcp in your firewall if your are testing externally. This does not apply to everyone.

## Run the API

At this point you should be able to deploy the application to the cluster and test the interface per the examples.

---
