# Pinflow Orbit

![Pinflow Orbit - In-memory geolocation database](/assets/banner.png)

A versatile, in-memory geolocation database designed for real-time tracking and management of courier location data within the Pinflow platform.

## Overview

**Pinflow Orbit** is a microservice written in Go that stores real-time geolocation data in-memory and provides a RESTful API for querying and managing the data.

### Main Objective

It is primarily designed to store live courier location data and provide a way for the main server to query and update the data using a RESTful API.

Orbit is designed to be a standalone service that can be deployed independently of the main Pinflow server.

To use the service, the main server needs to authenticate requests using a shared secret key (a simple API key is used for now).

### Tech Stack

Written in Go with the standard library. Built and run using Docker.

## Getting Started

### Prerequisites

- Go 1.22
- Docker (optional, for containerized deployment)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/katistix/pinflow-orbit.git
    cd pinflow-orbit
    ```

2. Set up environment variables:

    Create a `.env` file in the root of the project and add your API key:

    ```env
    API_KEY=your_api_key_here
    ```

### Running the Service

#### Locally

To run the service locally:

1. Ensure you have the required Go version installed.
2. Run the service:

    ```bash
    go run main.go
    ```

#### Using Docker

To build and run the service using Docker:

1. Build the Docker image:

    ```bash
    docker build -t pinflow-orbit .
    ```

2. Run the Docker container:

    ```bash
    docker run -d -p 8080:8080 --env-file .env pinflow-orbit
    ```

### API Endpoints

The following endpoints are available:

#### Health Check

- **GET /health**

    Checks if the service is running.

    ```sh
    curl -X GET http://localhost:8080/health
    ```

#### Authenticated Health Check

- **GET /auth-health**

    Checks if the service is running and requires authentication.

    ```sh
    curl -X GET http://localhost:8080/auth-health -H "Authorization: Bearer your_api_key_here"
    ```

#### Get All Locations

- **GET /locations**

    Retrieves all stored locations.

    ```sh
    curl -X GET http://localhost:8080/locations -H "Authorization: Bearer your_api_key_here"
    ```

#### Get Location

- **GET /location**

    Retrieves a location for a specific user.

    Parameters:
    - `userId`: The ID of the user whose location is being retrieved.

    ```sh
    curl -X GET "http://localhost:8080/location?userId=user1" -H "Authorization: Bearer your_api_key_here"
    ```

#### Set Location

- **POST /set**

    Sets or updates a location for a user.

    Request Body:
    ```json
    {
        "userId": "user1",
        "longitude": 10.0,
        "latitude": 20.0
    }
    ```

    ```sh
    curl -X POST http://localhost:8080/set -H "Authorization: Bearer your_api_key_here" -H "Content-Type: application/json" -d '{"userId": "user1", "longitude": 10.0, "latitude": 20.0}'
    ```

#### Delete Location

- **POST /delete**

    Deletes a location for a user.

    Request Body:
    ```json
    {
        "userId": "user1"
    }
    ```

    ```sh
    curl -X POST http://localhost:8080/delete -H "Authorization: Bearer your_api_key_here" -H "Content-Type: application/json" -d '{"userId": "user1"}'
    ```

### Testing

#### Running Unit Tests

To run the unit tests locally:

```bash
go test -v ./...
```

### Continuous Integration

A GitHub Actions workflow is set up to run tests on every commit to the main branch. The workflow configuration is as follows:

#### `.github/workflows/go.yml`

```yaml
name: Go

on:
  push:
    branches: 
      - master
  pull_request:
    branches: 
      - master

jobs:
  test:

    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.22

    - name: Get dependencies
      run: go mod download

    - name: Run tests
      run: go test -v ./...
```

### Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature-name`).
3. Make your changes and commit them (`git commit -am 'Add new feature'`).
4. Push to the branch (`git push origin feature/your-feature-name`).
5. Create a new Pull Request.

### License

This project is licensed under the MIT License.