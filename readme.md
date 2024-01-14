demo project web-based app file storage

the architecture is rather simple , it use postgres for store user account and minio for blob store  


use docker compose for test 
it will pull dependency image (postgres, minio, caddy) and also build image of the app from Dockerfile

```
    docker compose up -d
```

try "localhost:2080".

