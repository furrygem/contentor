# Contentor

furrygem/contentor is a WIP media storage / sharing service. Uses microservice architecture.

## Deployment

furrygem/contentor is deployed using docker-compose

### Configuration

#### Microservices

Refer to auth-service/README.md and content-service/readme.md for services configuration instructions

#### Nginx

Nginx configuration files can be found in nginx/ directory

### Running

```bash
git clone https://github.com/furrygem/contentor
cd contentor

# Configuration steps

docker-compose up
```

## What is running?

|      Service name     |          Image / Build          |                         Brief description                        |    Ports policy*    | Volumes* |       Secrets       |
|:---------------------:|:-------------------------------:|:----------------------------------------------------------------:|:-------------------:|:--------:|:-------------------:|
|          etcd         | _img_: docker.io/bitnami/etcd:3 | Etcd key-value storage for services communication during startup |      2379:2379      |     -    |          -          |
|    auth-service-db    |           _img_: mysql          |                    MySQL DBMS for auth-service                   |      3306:3306      |     -    | mysql-root-password |
|        adminer        |          _img_: adminer         |             Adminer instance for database inspection             |      8090:8080      |     -    |          -          |
|      auth-service     |       _bld_: auth-service       |       Auth service provides authenticaion and authorization      |      8000:8000      |     -    |          -          |
| content-service-minio |    _img_: quay.io/minio/minio   |                   Minio S3-like object storage                   | 9000:9000 9001:9001 |     -    |          -          |
|    content-service    |      _bld_: content-service     |           Provides API to store and fetch media content          |      8001:8000      |     -    |          -          |
|         nginx         |           _img_: nginx          |  Nginx web server providing reverse proxy for the microservices  |       8088:80       |     -    |          -          |

\* by default
