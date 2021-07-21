# How to run the K8s example with a local Docker Database

The K8s example for the course uses a hosted database for testing.

This README describes how to use the local Docker database used in the previous examples.

This example was tested on a Linux server running a local cluster. It has not been tested anywhere else. 

The steps fit into the course around the CONTAINERS AND KUBERNETES section after completing the "Defining our K8s Deployment" video.

---

### TLDR - 

The quick description is that you can use the existing docker-compose approach to start the database (from the earlier examples) and then modify the service to use a different port (so that it does not conflict with the docker-compose app container that is running).

### Deep Dive

This can be done in 3 easy steps.

1. Change the port in the `config/service.yml` file. 

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

2. Create the service.

    ```
    kubectl apply -f config/service.yml
    ```

    Note: If you already created the service per the tutorial then you will need to delete the original service using `docker delete service comments-api`.

3. Change the syntax of the `port-forward` command to point to the new port (8888).

    ```
    kubectl port-forward --address 0.0.0.0 service/comments-api 8888
    Forwarding from 0.0.0.0:8888 -> 8080
    ```

    Note: Here we are using an additionl argument `--address 0.0.0.0` to allow the connection to bind on all interfaces. This is to faciliate testing from external hosts and is not necessary if you are running everything locally. Note that you may need to open port 8888/tcp in your firewall if your are testing externally. This does not apply to everyone.

## Run the API

At this point you should be able to deploy the application to the cluster and test the interface per the examples substituting port `8080` with `8888`.

For example:

```
$ curl http://127.0.0.1:8888/api/health
{"Message":"I am Alive","Error":""}
```

Remember to change the port in your Postman URLs as well when testing against the cluster.

Another example:

```
$ curl -X GET http://127.0.0.1:8888/api/comment
[{"ID":1,"CreatedAt":"2021-07-18T16:06:03.944902Z","UpdatedAt":"2021-07-18T16:06:03.944902Z","DeletedAt":null,"Slug":"/hi","Body":"Some message 1.","Author":"Bill"},{"ID":2,"CreatedAt":"2021-07-18T16:06:13.97669Z","UpdatedAt":"2021-07-21T15:58:34.67692Z","DeletedAt":null,"Slug":"/testing-put-update","Body":"Some NEW message 2.","Author":"Bill"},{"ID":3,"CreatedAt":"2021-07-19T14:35:46.162889Z","UpdatedAt":"2021-07-21T15:23:29.121784Z","DeletedAt":null,"Slug":"/testing-put-update","Body":"Some NEW message 3.","Author":"Bill"}]
```

---
