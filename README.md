# Calculation Service
This service is meant to calculate the cost of a trip based on the ORIGIN and DESTINATION. We use [Gorilla Mux](https://github.com/gorilla/mux) as a framework for the REST API.

The service currently depends on the [gmaps-adapter](https://github.com/AITestingOrg/gmaps-adapter).

## Set Up

### Go Environment
Be sure to clone the repository (calculation-service) in the proper folder (in this case under the src folder). Your structure should look somewhat like so:
```go
go (root)
 │ - src
 │     │ - project_files (project root)
 │     └───
 │ - pkg
 │     │ - ...
 │     └───
 │ - bin
 │     │ - ...
 │     └───
 └───
```

More information on setting up the proper go environment can be found in the official documentation.

### Running 
To run the service **locally**, make sure to have the gmaps-adapter running and then continue with these steps:
  - In the root of the project, run `go get -u github.com/gorilla/mux`.
  - Afterwards run `go build`.
  - Then lastly run `./calculation-service`.

To run the service in **containers** using Docker (make sure you have Docker installed and running), clone the repository and simply run: `docker-compose up --build` at the root of the project.