# go-user-service

A gRPC-based user service built with Go, using MongoDB as database, supporting JWT with RSA keys.
## Overview

`go-user-service` is a lightweight microservice designed for user authentication and authorization. It is built with:

- **gRPC** for fast and efficient communication  
- **JWT (RSA256)** for secure token-based authentication  
- **MongoDB** as the data store for users, roles, and permissions  

---

## Features

- User, Role, Permission management via gRPC
- MongoDB backend
- JWT authentication with RSA keys (public/private)
- Dockerized for easy deployment
- Ready for Docker Hub usage and Docker Compose orchestration

---

## Getting Started

### Prerequisites

- Docker (if running via container) or Go installed locally  
- MongoDB instance running and accessible  
- RSA256 key pair (private and public keys) for JWT signing  

---

### Generate RSA Key Pair

Before running the service, you must generate an RSA key pair (private and public keys) to sign and verify JWT tokens.

You can generate the keys using OpenSSL:

```bash
# Generate 2048-bit RSA private key
openssl genrsa -out private.key 2048

# Generate the corresponding public key
openssl rsa -in private.key -pubout -out public.key
```

Place these key files in a secure directory on your machine, for example:

C:/Projects/my/go-user-service/keys/private.key  
C:/Projects/my/go-user-service/keys/public.key

Set Environment Variables

Set environment variables to let the service know where your keys and MongoDB URI are:
Windows PowerShell example
```bash
$env:PRIVATE_KEY_FILE = "C:/Projects/my/go-user-service/keys/private.key"
$env:PUBLIC_KEY_FILE = "C:/Projects/my/go-user-service/keys/public.key"
$env:MONGO_URI = "mongodb://localhost:27017"
$env:GRPC_PORT = "50051"
```
Linux / macOS Bash example
```bash
export PRIVATE_KEY_FILE="/path/to/your/private.key"
export PUBLIC_KEY_FILE="/path/to/your/public.key"
export MONGO_URI="mongodb://localhost:27017"
export GRPC_PORT="50051"
```
Pull and Run Docker Container

After keys are generated and environment variables are set, you can run the service using Docker:

```bash
docker pull yasinsaeeniya/go-user-service

docker run -p 50051:50051 \
  -e PRIVATE_KEY_FILE="C:/Projects/my/go-user-service/keys/private.key" \
  -e PUBLIC_KEY_FILE="C:/Projects/my/go-user-service/keys/public.key" \
  -e MONGO_URI="mongodb://host.docker.internal:27017" \
  -e GRPC_PORT="50051" \
  yasinsaeeniya/go-user-service
```
    Note:
    When running Docker on Windows or Mac, use host.docker.internal as the MongoDB host to connect to your local MongoDB instance.

Running Locally (without Docker)
