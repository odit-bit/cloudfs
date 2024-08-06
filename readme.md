 ## Getting Started
 Remote file storage application.
 
### Clone repository
the simplest way to try this, is directly build from the source or use `go run` command

```shell
#clone git
git clone https://github.com/odit-bit/cloudfs.git

# run command
cd cloudfs
go run ./cmd/web 
```
it will run the web server (html page) with all service using in-memory implementation and will be blank-state for every time it start [visit localhost:8181 ](http://localhost:8181) 

#### Feature
- User authentication
- Standard file operation (upload, download, delete).
- Share the file using time expiration generated url.
- include built-in implementation for testing

## Architecture
<img title="a title" alt="Alt text" src="cloudfs-simple-diagram.jpg">


### Environment Variable
the architecture has 3 core implementation for managing data.
- user (postgres)
- blob (minio)
- token (redis)  

http dependency:
- session-cookie (redis)

#### Blob
```shell
BLOB_MINIO_ENDPOINT
BLOB_MINIO_ACCESS_KEY
BLOB_MINIO_SECRET_ACCESS_KEY
```
#### User
```shell
USER_PG_URI
```

#### Token
```shell
TOKEN_REDIS_URI
```

#### Session
```shell
SESSION_REDIS_URI
```

### Docker compose
Example to use remote infrastructure for the service with caddy as reverse proxy, Docker will build from `Dockerfile` and fetch the neccessary image. [visit localhost:2080 ](http://localhost:2080) 
```shell
# start
docker compose up -d

# stop
docker compose down -v
```

