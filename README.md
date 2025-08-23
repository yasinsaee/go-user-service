# 🚀 go-user-service

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

```bash
curl -O https://raw.githubusercontent.com/yasinsaee/go-user-service/main/docker-compose.yml
```

### Step 3: Start the services
```bash
docker-compose up -d
```

