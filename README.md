# pinflow-orbit

A versatile, in-memory geolocation database designed for real-time tracking and management of courier location data within the Pinflow platform.

## Overview

**Pinflow Orbit** is a microservice written in Go that stores realtime geolocation data in-memory and provides a RESTful API for querying and managing the data.

### Main Objective

It is primarily designed to store the live courier location data and provide a way for the main server to query and update the data using a RESTful API.

Orbit is designed to be a standalone service that can be deployed independently of the main Pinflow server.

To be able to use the service, the main server needs to authenticate the requests using a shared secret key. (A simple API key is used for now)

### Tech Stack

Written in Go with the standard library. And built and run using Docker.


## Building the Docker Image
```bash
docker build -t pinflow-orbit:[version_tag] .
```

## Running the Docker Container
```bash
docker run -p 8080:8080 --env-file .env.prod pinflow-orbit:[version_tag]
```
> Note: The environment variables are loaded from a `.env.prod` file. Make sure to create this file and pass it to the container.



## RESTful API

The API is designed to be simple and easy to use. It provides endpoints for querying and updating the courier location data.

### Endpoints

#### `GET /location?userId=<userId>`
- Returns the latest location data for the given user.
- Requires authentication using the shared secret key.
- Returns a 404 if the user is not found.

**Response example:**
```json5
{
  "latitude": 12.345,
  "longitude": 67.890,
  "last_update": 1719736810 // Unix timestamp
}
```

#### `POST /set`

- Updates the location data for the given user.
- Requires authentication using the shared secret key.
- The request body should contain the following JSON data:
  ```json5
  {
    "userId": "123",
    "latitude": 12.345,
    "longitude": 67.890
  }
  ```
  
- Returns a 200 if the data is updated successfully.
- Returns a 400 if the request body is invalid.
- Returns a 401 if the authentication fails.

**Response example:**
```json5
{
    "message": "Location updated"
}
```

#### `GET /locations`

- Returns the latest location data for all users.
- Requires authentication using the shared secret key.
- The response is a JSON object where the keys are the user IDs and the values are the location data.

**Response example:**
```json5
{
  "user_id_123": {
    "latitude": 12.345,
    "longitude": 67.890,
    "last_update": 1719736810
  },
  "user_id_456": {
    "latitude": 12.345,
    "longitude": 67.890,
    "last_update": 1719736810
  }
}
```