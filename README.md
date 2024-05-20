[![codecov](https://codecov.io/gh/Neniel/gotennis/graph/badge.svg?token=DOLX5B94B0)](https://codecov.io/gh/Neniel/gotennis)


# Running the application on Kubernetes
### 1. Install minikube
### 2. Start minikube
Run the command below to start minikube
```
minikube start
```

### 2. Configure minikube
Run the command below in order to configure your minikube:
```
make k8s-configure
```

### 3. Create a `secrets.yml` file in the `kubernetes` directory:
```yml
apiVersion: v1
kind: Secret
metadata:
  name: application-secret
type: Opaque
data:
  mongodb_connection_string: <BASE64_MONGODB_CONNECTION_STRING>
  redis_server_address: <BASE64_REDIS_SERVER_ADDRESS>
  redis_password: <BASE64_REDIS_PASSWORD>
  grafana_graphite_token: <BASE64_GRAFANA_GRAPHITE_TOKEN>
```

### 4. Configure /etc/hosts

Add the following line to the end of your `/etc/hosts` file
```
127.0.0.1 categories.tennis.dev players.tennis.dev
```

### 5. Deploy on your Kubernetes cluster!
Run the command below to start the deployment:
```
make k8s-deploy
```

### 6. Add some magic
Open a new terminal and run the command below to add some magic
```
minikube tunnel
```

# Running the application with docker-compose

Run the command below to start the deployment using docker-compose:

```
make deploy
```
