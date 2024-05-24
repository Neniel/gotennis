[![codecov](https://codecov.io/gh/Neniel/gotennis/graph/badge.svg?token=DOLX5B94B0)](https://codecov.io/gh/Neniel/gotennis)

# Running the application on Kubernetes
### 1. Install minikube or kind
### 2. Start a k8s cluster
Run the command below to create a k8s cluster using minikube
```
make k8s-configure-minikube
```

Run the command below to create a k8s cluster using kind
```
make k8s-configure-kind
```

### 2. Configure your k8s cluster
Run the command below in order to configure your minikube k8s cluster:
```
make k8s-configure-minikube
```

Run the command below in order to configure your kind k8s cluster:
```
make k8s-configure-kind
```

### 3. Configure the Heml Chart
Create a `value-dev.yml` file in the `.deployment/k8s/tennis-app` out of `values-blueprint.yaml`


### 4. Configure /etc/hosts

Add the following line to the end of your `/etc/hosts` file
```
127.0.0.1 categories.tennis.dev players.tennis.dev
```

### 5. Deploy the app on your Kubernetes cluster!
Run the command below to start the deployment:
```
make k8s-deploy
```

### 6. Additional step if you're using minikube
Open a new terminal and run the command below to add some magic
```
minikube tunnel
```

# Running the application with docker-compose

Run the command below to start the deployment using docker-compose:

```
make deploy
```
