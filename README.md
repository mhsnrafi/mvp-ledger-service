## Challenge Bitburst MVP Ledger Service
### Overview
The Bitbrust Ledger Service is a system designed to manage user balances, transactions, and check transaction history. The system is built using Golang and leverages Docker and Docker Compose for containerization and orchestration. The database used for storing user and transaction information is PostgreSQL, while Redis serves as a caching layer for the application.

## System Components
The system consists of the following components:

1. **Ledger Service**: A Golang application that exposes RESTful API endpoints for managing Funds, transactions history, and account balances. 
2. **PostgreSQL**: A relational database for storing user and transaction data. 
3. **Redis**: An in-memory data store for caching and improving the performance of the system. 
4. **Nginx**: A reverse proxy that directs incoming traffic to the appropriate ledger service accross multiple instances. 
5. **End-to-End Tests**: A testing suite to validate the functionality and performance of the ledger service.
6. **Docker**: The system is containerized using Docker and Docker Compose, making it easy to build, deploy, and scale the application in any environment.

## Functional Requirements
1. Allow users to add funds to their account
2. Display a user's transaction history with pagination
3. Rate limit user requests to prevent abuse or excessive usage
4. Perform end-to-end testing to ensure the system behaves as expected 
5. Caller cannot guarantee that it will call exactly once for the same money transfer

## Non-Functional Requirements
1. High availability
2. Scalability
3. Fault tolerance
4. Maintainability

## Project Setup
1. Ensure you have the following software installed on your machine: 
   1. Docker: https://docs.docker.com/get-docker/
   2. Docker Compose: https://docs.docker.com/compose/install/
   3. Golang: https://golang.org/doc/install
2. Clone the project repository and navigate to the project directory:
```git
git clone <repository-url>
cd <project-directory>
 ```

## Running the Project
1. Build and start the project using Docker Compose:
```dockerfile
docker-compose up --build -d --scale ledger-service=3
```
This command will build and start the following services:

* 3 instances of the ledger service
* PostgreSQL
* Redis
* Nginx
* Run End to End integration tests and stop the test container

2. To check the status of the running containers, run:
```dockerfile
docker-compose ps
```
3. Access the ledger service API through the Nginx reverse proxy by sending requests to **http://localhost:4000**


## Assumptions
The following assumptions have been made while designing the funds balance and transaction history service:

1. Users are identified by a unique ID (UID) provided in the API requests.
2. The system is expected to handle thousands of concurrent users with hundreds of transactions each.
3. Transaction history data is frequently accessed and benefits from caching. 
4. The service is primarily focused on handling user funds and transaction history, with no additional features or requirements. 
5. The primary data storage is a relational database, and Redis is used for caching purposes.
7. The service will be deployed on multiple instances behind a load balancer for high availability and fault tolerance with the help of nginx



## Component Interaction or Flow

1. Clients send HTTP requests to the Nginx reverse proxy.
2. Nginx directs incoming requests to one of the ledger service instances. 
3. A unique transaction ID is generated 
4. Funds are added to the user's account using a distributed lock to ensure consistency
5. The ledger service interacts with the PostgreSQL database to store, update, and retrieve user and transaction data. 
6. The ledger service uses Redis to cache frequently accessed data for improved performance. 
7. End-to-end tests simulate client requests to the ledger service to validate the system's functionality
8. When the Docker container starts, it will automatically execute the end-to-end test cases and then stop the container.

### For retrieving transaction history:
1. Data is fetched from the Redis cache if available 
2. If not available in the cache, data is fetched from the database 
3. Pagination is applied, and the response is returned to the client

## Deployment
The service can be deployed on multiple instances behind a load balancer to ensure high availability and fault tolerance. Horizontal scaling can be used to handle increased load.

## Logging
The service included logging to track errors. This will help identify bottlenecks, diagnose issues, and improve the overall system.

## Security


**Note: JWT implementation is currently commented out in the project. Uncommenting the code will enable JWT-based authentication for the APIs.**

To make the APIs secure, we can use JSON Web Tokens (JWT) for authentication. The following approach can be taken to generate and refresh access tokens:

* The GenerateAccessTokens function creates two types of tokens - an access token and a refresh token - for a given email.
* The access token has a set expiration time, after which it will no longer be valid. This expiration time can be configured by setting the JWT_ACCESS_TOKEN_EXPIRATION_TIME environment variable.
* The refresh token also has a set expiration time, after which it will no longer be valid. This expiration time can be configured by setting the JWT_REFRESH_TOKEN_EXPIRATION_TIME environment variable.
* The function calls the CreateToken function twice to create both the access and refresh tokens and returns them.
* If there is an error during the creation of either token, the function returns an error. 
* By using JWT-based authentication, we can secure our APIs and ensure that only authorized users can access them.

## API Documentation
To test the API endpoints directly from the documentation, making it easier to ensure that the API is working as expected build swagger api documentationa  user-friendly interface to quickly understand the API’s capabilities and functions
```bash
http://localhost:4000/swagger/index.html#/
```

#### Postman API Collection is added in the project directory - MVP ledger Service Api Collection.postman_collection.json

