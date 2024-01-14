demo project web-based app file storage

the architecture is rather simple , it use postgres for store user account and minio for blob store  

## Getting Started
use docker compose for test 
it will pull dependency image (postgres, minio, caddy) and also build image of the app from Dockerfile

```
    docker compose up -d
```

try "localhost:2080" in browser.  
to manage the uploaded file it can access from "localhost:9090", it is the minio admin dashboard.  
login with user "admin" and "admin12345", for change the root user and password just change the compose file before. 
[minio-website](https://min.io)