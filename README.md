# Happy birthday application
This project runs a containerized Golang Application with the [Gin Web Framework](https://github.com/gin-gonic/gin).

## How it works
The application stores and prints the User's birthday, based on the endpoint and http method that you send. The app is listening on `http://localhost:8080`

### Validations
Before storing a new user into the datbase, the application validates:

- If the username already exist into the DB
- If the username only contain leters
- If the date of Birth is in the correct format `YYYY-MM-DD`
- If the day of birth is not greater than today

### Deployment
This procedure intend to deploy the Kubernetes resources `Ingress`, `Service` and `Deployment` with the help of [helm](https://helm.sh/), into a running k8s cluster on AWS. With this procedure, we can guarantee a non downtime deployment.  

The deployment consists into two parts, first testing, building and publishing the new docker image to the Docker image repository.

The second consist on deploying the kubernetes resources into the cluster.

## Cloud diagram
![Cloud diagram](hbday-cloud.png "Cloud diagram")

## Cheat Sheet
```shell
# Run locally with docker compose
# This command spins up two docker containers
# a postgres db and the application
$ make run

# Test your applicaiton locally
$ make test 

# Test, Build and Push your image to the Docker Hub
$ make push

# Test, Build and Push  your docker image to Docker hub and Deploy your k8s resopurces into the cluster
$ make deploy

# Create new user
$ curl http://localhost:8080/hello/<username> --include \
  --header "Content-Type: application/json"  \
  --request "PUT" \ 
  --data '{"date_of_birth": "YYYY-MM-DD"}'`

# Check existing user and how long for its birthday
$ curl http://localhost:8080/hello/<username>
```

## Prerequists and assumptions
### Run locally
- [Docker](https://docs.docker.com/get-docker/) installed in your local machine
- [kubectl](https://kubernetes.io/docs/tasks/tools/) installed in your local machine
- [helm](https://helm.sh/docs/intro/install/) Installed on your local machine
- [go](https://golang.org/doc/install) Installed on your local machine
  
### Run in the cloud (all above requirements, plus the below ones)
- An already running k8s cluster running on AWS 
- [Nginx Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/) already installed on your k8s cluster
- A proper host configured in front of NLB 
- A proper host configured on Ingress (Update the `ingress.host` in [backend-chart/values.yaml](backend-chart/values.yaml))
- Postgres database already running with DB `hbday_db` (Update the `env` in [backend-chart/values.yaml](backend-chart/values.yaml) in order to set correct db configuration)
