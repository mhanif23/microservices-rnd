# Library Apps — Playbook v1.0 Boilerplate

- user-svc (Go + Postgres)
- book-svc (NestJS + Mongo)
- borrow-svc (Spring Boot + Postgres)

## Quick Start (Docker Compose)
```bash
cd ops
docker compose up --build
# user: http://localhost:3000/health
# book: http://localhost:3001/health
# borrow: http://localhost:3002/health
# Adminer: http://localhost:8080  | Postgres host: postgres, user/pass: library/library, db: librarydb
# Mongo: localhost:27018
```

Seed cepat (Adminer → SQL):
```sql
INSERT INTO users(id,email,name,membership_status) VALUES
('11111111-1111-1111-1111-111111111111','demo@library.io','Demo User','ACTIVE');

INSERT INTO inventory(book_id,total,available) VALUES
('22222222-2222-2222-2222-222222222222',5,5);
```

Coba endpoints:
```bash
# User
curl localhost:3000/health
curl -X POST localhost:3000/users -H "content-type: application/json" -d '{"email":"a@b.c","name":"Alice"}'

# Book (ingat gunakan _id UUID string saat create)
curl localhost:3001/health
curl -X POST localhost:3001/books -H "content-type: application/json" -d '{"_id":"22222222-2222-2222-2222-222222222222","title":"Microservices 101"}'

# Borrow
curl localhost:3002/health
curl -X POST localhost:3002/borrows -H "content-type: application/json" -H "Idempotency-Key: demo-1"   -d '{"user_id":"11111111-1111-1111-1111-111111111111","book_id":"22222222-2222-2222-2222-222222222222"}'
```

## Minikube (K8s) — ringkas
```bash
kubectl apply -f deploy/k8s/namespace.yaml

# DB via Helm (Bitnami) - jalankan ini dulu
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install postgresql bitnami/postgresql -n lib-dev   --set auth.username=library,auth.password=library,auth.database=librarydb

helm install mongodb bitnami/mongodb -n lib-dev --set auth.enabled=false

# Build image ke daemon minikube
eval $(minikube docker-env)
docker build -t user-svc:0.1.0 services/user-svc-go
docker build -t book-svc:0.1.0 services/book-svc-ts
docker build -t borrow-svc:0.1.0 services/borrow-svc-spring

# Deploy apps
kubectl apply -f deploy/k8s/user/deployment.yaml
kubectl apply -f deploy/k8s/book/deployment.yaml
kubectl apply -f deploy/k8s/borrow/deployment.yaml

# Ingress (Traefik addon harus enable)
kubectl apply -f deploy/k8s/traefik/ingress.yaml
```

## Contracts
Lihat folder `/contracts` (OpenAPI v1). Disarankan generate client stubs dari sini.
