# Running the example with a local DB

The K8s examples for the course used a hosted database for testing.

This file describes how to use a local Docker database to do the same. 

This example was tested on a Linux server running a local cluster. It has not been tested anywhere else. 

These steps pickup in the CONTAINERS AND KUBERNETES section of the course after completing the Defining our K8s Deployment.


## Prerequisite Steps

1. If you are still running the docker-compose example - then you will need to stop it. (If you don't then you may run into issues with conflicting ports.) To stop it type: '<Ctrl><c>' in the window where it is running. 

2. You may need to clean up your docker environment to remove stopped containers. (If you don't you may run into issues with conflicting container names.) To see the containers you can use: `docker ps -a`. Look for containers with the name "comments-api" and remove them by hand by using the `docker rm` command and their container id. For example: 

    ```
    [bgercken@sparta go-rest-api]$ docker ps -a
    CONTAINER ID        IMAGE                  COMMAND                  CREATED             STATUS                     PORTS               NAMES
    3b20777f999a        go-rest-api_api        "./app"                  16 seconds ago      Exited (2) 4 seconds ago                       comments-restapi
    48a1b9c2c72a        postgres:12.2-alpine   "docker-entrypoint.s…"   17 seconds ago      Exited (0) 3 seconds ago                       comments-database
    da0ccf1893ed        3daa8f4acc31           "/bin/sh -c 'CGO_ENA…"   4 hours ago         Exited (2) 4 hours ago                         sad_sinoussi
    3b1d40967899        ec1613ae3afc           "/bin/sh -c 'CGO_ENA…"   4 hours ago         Exited (2) 4 hours ago                         hopeful_tharp
    21e902253b60        3724ea0cf264           "/bin/sh -c 'CGO_ENA…"   4 hours ago         Exited (2) 4 hours ago                         exciting_ellis
    1f35cb39ff43        10960de6cb63           "container-entrypoin…"   5 days ago          Exited (130) 5 days ago                        zealous_brahmagupta
    0296d8032975        98c4c34c4dd2           "/bin/sh -c 'yum ins…"   2 weeks ago         Exited (1) 2 weeks ago                         nice_johnson
    cbe0a87ad7f3        8edf6234b95c           "/bin/sh -c '/usr/lo…"   2 weeks ago         Exited (2) 2 weeks ago                         charming_shamir

    [bgercken@sparta go-rest-api]$ docker rm 3b20777f999a 48a1b9c2c72a
    3b20777f999a
    48a1b9c2c72a
    ```

3. For this example I opened TCP port 8080 in my firewall to allow connections from the outside.

    On Linux you can use the following (depending on your version).

    As root:
    ```
    firewall-cmd --permanent --add-port=8080/tcp
    firewall-cmd --reload
    ```

4. We need to revert the `sslmode=required` change (shown around minute 7:15 of the video on `Defining our K8s Deployment`. Change: `sslmode to disable`.

5. Rebuild your container and tag it for your docker repository. `docker build -t mydockerhubname:/comments-api:latest .` 

6. Push the container to your repository: `docker push mydockerhubname:/comments-api:latest`.

## DB Setup and Statup Steps:

1. Create a script for starting the database. We are assuming this is a fresh start - so a docker volume will be created to store the database the first time that the script is executed. The database will be created the first time you deploy the pods.

    ```
    cat > run-db.sh << EOF

    #!/usr/bin/env bash
    #
    . ./ENV.sh
    docker volume create comments-api-db_postgres

    docker run --name comments-api-db \
      -v comments-api-db_postgres:/var/lib/postgresql/data \
      -e POSTGRES_PASSWORD=$DB_PASSWORD -p $DB_PORT:$DB_PORT -d postgres:12.2-alpine 
    EOF
    ```

2. Create an environment file (as shown in the course substituting your choice for values) and your IP address.

    ```
    cat > ENV.sh << EOF
    export DB_USERNAME=postgres
    export DB_PASSWORD=postgres
    export DB_HOST=192.168.0.100
    export DB_TABLE=postgres
    export DB_PORT=5432
    EOF
    ```

3. Change the permissions on the database script to make it executable. `chmod +x run_db.sh`


4. Start the database. 

    ```
    ./run_db.sh
    ```

## Now follow the steps to deploy the application

1. Create a script to run the application:

    ```
    cat > deploy_app.sh << EOF
    #!/usr/bin/env bash

    . ./ENV.sh
    envsubst < config/deployment.yaml | kubectl apply -f -
    EOF

2. Change the permissions on the deploy script to make it executable. `chmod +x deploy_app.sh`

3. Start the application.

    ```
    ./deploy_app.sh
    ```

4. Forward the port as in the example but also use the `--address` argument so that you can access the API at the local IP address (rather than local host). See note below.

    ```
    kubectl port-forward --address 0.0.0.0 service/comments-api 8080:8080
    ```

Note: In step 4 you are exposing the application to hosts outside of the local host. I did this to facilitate testing from within my network (so that I can use postman). Be aware that other people can access your API and do not do this in production unless it is something that you specifically want to do. (You probably won't be using `kubectl port-forward ...` in production anyway. :-)



