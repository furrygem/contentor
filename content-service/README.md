# Content Service

Serves content using MinIO

## API Reference

| Path              | Method | Description                                         |
|-------------------|--------|-----------------------------------------------------|
| /api/objects      | GET    | Get a list of objects in the storage                |
| /api/objects      | POST   | Upload an object to the storage via FormData `file` |
| /api/objects/{id} | GET    | Download a file by its key                          |
| /api/objects/{id} | DELETE | Delete a file by its key                            |

### Authorization

#### Token

JWT Token, signed with EdDSA algorithm, containing `sub` field containing the information about the resource owner

```json
{
    "typ": "jwt",
    "alg": "EdDSA"
}
{
    "sub": "user_id"
}
```

## Startup script

**generate** - fully generates private and public keys on startup

**use-private** - use specified private key to derive public key on startup

**use-http** - fetch public key from an HTTP endpoint

**use-etcd** -  fetch public key form etcd key-value store

For more information run

```bash
./start.py -h
```

or

```bash
start.py [script] -h
```

## Configuration

Configuration is done using environment variables or .env file

| Environment variable                       | Description                                                      | Default        | Required |
|--------------------------------------------|------------------------------------------------------------------|----------------|----------|
| CONTENT_SERVICE_AUTH_EDDSA_PUBLIC_KEY_FILE | Public EdDSA key file to be used for JWT signature verification  | key.pem        | False    |
| CONTENT_SERVICE_MINIO_WORKERS              | Amount of goroutines to work on file upload                      | 1              | False    |
| CONTENT_SERVICE_MINIO_CHANNEL_CAPACITY     | Capacity of the file upload channel                              | 1              | False    |
| CONTENT_SERVICE_MINIO_BUCKET_NAME          | Name of minio bucket to use for file upload                      | files          | False    |
| CONTENT_SERVICE_MINIO_ENDPOINT             | Minio server hostname                                            | minio          | False    |
| CONTENT_SERVICE_MINIO_ACCESS_KEY           | Minio access key                                                 |                | True     |
| CONTENT_SERVICE_MINIO_SECRET_KEY           | Minio secret key                                                 |                | True     |
| CONTENT_SERVICE_MINIO_USE_SSL              | Use of SSL for minio connection                                  | False          | False    |
| CONTENT_SERVICE_SERVER_LISTEN_ADDR         | Server listening address                                         | 127.0.0.1:8080 | False    |