# Auth-Service

Authentication / Authorization provider for Contentor system

## Configuration

Add mysql root password

```bash
echo "superSecretPassword~" > mysql-root-password.txt
```

Populate `.env` [dot env] file

```
AUTH_SERVICE_DB_HOST=db
AUTH_SERVICE_DB_USER=root
AUTH_SERVICE_DB_PASS=insecure
AUTH_SERVICE_DB_NAME=auth-service
```

## Running

To run the service run:

```bash
docker-compose up
```

## API Reference

The service uses FastAPI which generates OpenAPI specification, you can find them at `/docs` on the running auth-service server

## start.py

**-m** - Run alembic migrations on startup

**--generate** - Generate private key and overwrite teh existing one (if exists)

**-u** - Accumulate uvicorn options

### etcd

**--use-etcd** - Use etd to upload public key

**--etcd-host** - Etcd host. Default "localhost"

**--etcd-port** - Etcd port. Default 2379

**--etcd-key** - Etcd key to put the public key on. Default "contentor/public-key.pem"
