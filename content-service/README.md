# Content Service

Serves content using MinIO

## API Reference

| Path              | Method | Description                                         |
|-------------------|--------|-----------------------------------------------------|
| /api/objects      | GET    | Get a list of objects in the storage                |
| /api/objects      | POST   | Upload an object to the storage via FormData `file` |
| /api/objects/{id} | GET    | Download a file by its key                          |