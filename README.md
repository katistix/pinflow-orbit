# pinflow-orbit

A versatile, in-memory geolocation database designed for real-time tracking and management of courier location data within the Pinflow platform.

## Overview

**Pinflow Orbit** is a microservice written in Go that stores realtime geolocation data in-memory and provides a RESTful API for querying and managing the data.

### Main Objective

It is primarily designed to store the live courier location data and provide a way for the main server to query and update the data using a RESTful API.

Orbit is designed to be a standalone service that can be deployed independently of the main Pinflow server.

To be able to use the service, the main server needs to authenticate the requests using a shared secret key. (A simple API key is used for now)

### Tech Stack

Written in Go with the standard library.

## API

> Documentation is a work in progress.