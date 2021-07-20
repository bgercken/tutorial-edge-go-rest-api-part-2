# How to run the K8s example with a local Docker Database

The K8s example for the course uses a hosted database for testing.

This README describes how to use a local Docker database to do the same. 

This example was tested on a Linux server running a local cluster. It has not been tested anywhere else. 

The steps fit into the course around the CONTAINERS AND KUBERNETES section after completing the "Defining our K8s Deployment" video.


## Prerequisite Steps

1. If you are still running the docker-compose example - then you will need to stop it. (If you don't then you may run into issues with conflicting ports.) To stop it type: `<Ctrl><c>` in the window where it is running. 

2. You may need to clean up your docker environment to remove stopped containers. (If you don't you may run into issues with conflicting container names.)

    To see the stopped containers you can use: `docker ps -a`.

    Look for containers with names that contain "comments" and remove them by hand by using the `docker rm` command and the containers ID. For example: 

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

3. For testing I had to open port 8080/tcp in my firewall to allow connections from the outside. I did this to facilitate remote testing.

    Firewalls vary by Linux distribution and version. For my OS I used the following commands.

    As root:
    ```
    firewall-cmd --permanent --add-port=8080/tcp
    firewall-cmd --reload
    ```

4. Next, we need to revert the `sslmode=required` change made to the `internal/database/databases.go` connectionString. 
    (This was shown around minute 7:15 of the video on `Defining our K8s Deployment`.)

    Change: `sslmode=require` back to `sslmode=disable`.

5. Rebuild your container and tag it for your docker repository. Example: `docker build -t myname/comments-api:latest .` 

6. Finally, push the container to your repository. Example: `docker push myname/comments-api:latest`.

## Database Setup and Startup Steps:

1. Create an environment file (as shown in the course substituting values of your choice and your IP address).

    Example:

    ```
    cat > ENV.sh << EOF
    export DB_USERNAME=postgres
    export DB_PASSWORD=postgres
    export DB_HOST=192.168.0.100
    export DB_TABLE=postgres
    export DB_PORT=5432
    EOF
    ```

2. Create a script for starting the database. 

    We are assuming this is a fresh start - so a docker volume will be created to store the database the first time that the script is executed.

    The database will be created the first time you deploy the pods.

    Example:

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

3. Update the permissions on the database script to make it executable. Example: `chmod +x run_db.sh`


4. Start the database. 

    Example:

    ```
    ./run_db.sh
    ```

## Application Startup Steps

1. Create a script to run the application.

    Example:

    ```
    cat > deploy_app.sh << EOF
    #!/usr/bin/env bash

    . ./ENV.sh
    envsubst < config/deployment.yaml | kubectl apply -f -
    EOF

2. Update the permissions on the deploy script to make it executable. `chmod +x deploy_app.sh`

3. Start the application.

    ```
    ./deploy_app.sh
    ```

4. Forward the port as in the example but also use the `--address` argument so that you can access the API at the local IP address (rather than local host). See note below.

    ```
    kubectl port-forward --address 0.0.0.0 service/comments-api 8080:8080
    ```

## Post Testing Steps

After you are done testing your application you will want to stop the postgres database so that it will not conflict with the docker-compose testing in up comming lessons. (Logging with Logrus uses docker-compose.)

1. Type: `docker stop comments-api-db`.

2. Type: `docker rm comments-api-db`.

This will remove the container create by the `run-db.sh` script.


**Note:** 

1. In step 4 you are exposing the application to hosts outside of the local host. I did this to facilitate testing from within my network (so that I can use postman). Be aware that other people can access your API and do not do this in production unless it is something that you specifically want to do. (You probably won't be using `kubectl port-forward ...` in production anyway. :-)




