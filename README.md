# event-sourcing

## Overview
This is a small golang API to try out event sourcing! This repo contains a dockerfile where the eventstore gets created. The current image that is being used is an image that support ARM chips. If you don't have an ARM chip you can find an suiting image [here](https://github.com/eventstore/EventStore/pkgs/container/eventstore).

### What does it do
The application contains 3 endpoints, to create a new product stuck, update the stock and get the latest information about your product. It's only possible to update the amount of the stock. While deleting the API checks we enough stock is available. 

### Where to find the eventstore
When you run the eventstore docker you can find the admin portal [here](http://localhost:2113/)

## Prerequisites
- Go version 1.21.6
- Docker

## Installation
1. Clone the repository: `git clone https://github.com/jeffvo/event-sourcing.git`
2. Navigate to the project directory: `cd event-sourcing`
3. Install the dependencies: `go mod download`

## Running the Application
1. To start the application, run: `go run cmd/main.go`

## Built With
- [Go](https://golang.org/) - The programming language used.
- [Docker](https://www.docker.com/) - Used for containerization.

## Improvements
Although the API works certain improvements can be made
* Resolve potential race condition
* Swagger
* Authentication
* Create appsettings
