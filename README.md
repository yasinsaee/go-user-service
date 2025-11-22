 ____ ___  ______________________________    __________________________________   ____.____________ ___________
|    |   \/   _____/\_   _____/\______   \  /   _____/\_   _____/\______   \   \ /   /|   \_   ___ \\_   _____/
|    |   /\_____  \  |    __)_  |       _/  \_____  \  |    __)_  |       _/\   Y   / |   /    \  \/ |    __)_ 
|    |  / /        \ |        \ |    |   \  /        \ |        \ |    |   \ \     /  |   \     \____|        \
|______/ /_______  //_______  / |____|_  / /_______  //_______  / |____|_  /  \___/   |___|\______  /_______  /
                 \/         \/         \/          \/         \/         \/                       \/        \/ 

A gRPC-based **User Service** built with Go, using MongoDB as database, supporting JWT with RSA keys, fully dockerized and ready for Docker Hub & Docker Compose usage.

### Features

- User, Role, Permission management via gRPC  
- MongoDB backend  
- JWT authentication with RSA keys (public/private)  
- Dockerized for easy deployment  
- Ready for Docker Hub usage and Docker Compose orchestration  

---

### 📦 Quick Start

**Prerequisites:**

- [Docker](https://www.docker.com/get-started)  
- [Docker Compose](https://docs.docker.com/compose/install/)  

### Step 1: Pull the latest Docker image

```bash
docker pull yasinsaeeniya/go-user-service:latest
```

### Step 2: Download `docker-compose.yml`

The repository already contains a ready-to-use `docker-compose.yml` file. You can download it directly from GitHub:

#### Ubuntu / MacOS
```bash
curl -O https://raw.githubusercontent.com/yasinsaee/go-user-service/master/docker-compose.yml
```
#### Windows
```bash
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/yasinsaee/go-user-service/master/docker-compose.yml" -OutFile "docker-compose.yml"
```

### Step 3: Start the services
```bash
docker-compose up -d
```

### Step 4: Verify services

Check if the containers are running:

```bash
docker-compose ps
```

---

### 🎉 Congratulations  

The **go-user-service** is now up and running on your system! 🚀  
You can start sending **gRPC requests** to it and integrate it into your applications.  


### 🔗 Useful Links  

- [gRPC Quick Start](https://grpc.io/docs/languages/go/quickstart/)  
- [MongoDB Documentation](https://www.mongodb.com/docs/)  
- [Docker Hub – go-user-service](https://hub.docker.com/r/yasinsaeeniya/go-user-service)  
- [Docker Compose Documentation](https://docs.docker.com/compose/)  




