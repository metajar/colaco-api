# ColaCo-API

### This is temporarily hosted at https://cola.henrynetworks.com

### Docs are https://cola.henrynetworks.com/docs

## Overview
Welcome to the ColaCo-API, a versatile and secure API framework designed to manage vending machine operations. This project uses Go, Docker for deployment, JWT for secure authentication, and offers an interactive API documentation via ReDoc. It's built to be extensible, allowing for easy integration of various storage backends and customization of its authentication mechanisms.

# Key Features
- Flexible Storage Interface: Use any storage backend by implementing the VendingStorageInterface.
- JWT Authentication: Secure your API endpoints with JSON Web Tokens.
- Docker Deployment: Easily build and deploy with Docker.
- Interactive API Documentation: Access detailed API documentation through ReDoc.


### Building and Launching with Docker

1. **Build the Docker image:**
   ```
   docker build -t colaco-api .
   ```
2. **Run the server:**
   ```
   docker run -p 8080:8080 colaco-api
   ```


### Accessing the API Documentation
- Navigate to `http://localhost:8080/docs` to view the ReDoc documentation and explore available endpoints.

## Using the API
Before making requests, ensure you're authenticated (where required) by obtaining a JWT token.
You can leverage the Client in **Additional Tools** to use the `get-token` option to get a token.
This token can the be used to communicate to all the endpoints.

#### build

```bash
docker build -t cola .
```

#### run

```bash
docker run -p 8080:8080 cola
```

#### Accessing ReDoc Documentation API Server

You can access redoc when ran locally by going to:

http://127.0.0.1:8080/docs


---

## Additional Tools

For interacting with the ColaCo-API, we've provided a command-line client tool located in the `cmd/client` directory. This client simplifies testing and interacting with the API from the command line.

To access and learn more about the client tool, visit: [cmd/client](https://github.com/metajar/colaco-api/tree/main/cmd/client).

There is also in the folder included a postman_collection to interact with the API as well.

