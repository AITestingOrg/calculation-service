# Calculation Service
This service is meant to calculate the cost of a trip based on the ORIGIN and DESTINATION.

The service currently depends on the [gmaps-adapter](https://github.com/AITestingOrg/gmaps-adapter).

## Set Up

### Go Environment


More information on setting up the proper go environment can be found in the Golang documentation.

### Running 
To run the service **locally**, make sure to have the gmaps-adapter running and then continue with these steps:
  - In the root of the project, run `go get -u github.com/gorilla/mux`


To run the service in **containers**, clone the repository and simply run: `docker-compose up --build` at the root of the project.