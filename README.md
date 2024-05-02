# Salesforge assignment solution

## Prerequisites
- Docker

## Installation
1. Clone the repository: `git clone https://github.com/cybre/salesforge-assignment.git`
2. Change into the project directory: `cd salesforge-assignment`
3. Set up the database password: `echo "<some-password>" > db_password.txt` 

## Usage
1. Build the Docker image and run API and database via docker-compose
```bash
make run
```

2. Access the API at the default port at `http://localhost:3000`
3. Import the OpenAPI v3 spec into your API testing app of choice: `swagger.yaml`

## Additional Commands
- To stop the containers: `make stop`
- To clean up the project (remove containers, networks, and volumes): `make clean`
- To run unit tests: `make test`