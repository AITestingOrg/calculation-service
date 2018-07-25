# Calculation Service
This service is meant to calculate the cost of a trip based on the ORIGIN and DESTINATION. We use [Gorilla Mux](https://github.com/gorilla/mux) as a framework for the REST API.

The service currently depends on the [gmaps-adapter](https://github.com/AITestingOrg/gmaps-adapter).

## Set Up

### Go Environment

Please read carefully through the [Go documentation](https://golang.org/doc/install) on how to set up the main environment. After downloading the package, a site should pop up giving you further installation/setup instructions.

To get the calculation service project (make sure you've set up the Go environment) simply cd into your GOPATH and run this command: `go get github.com/AITestingOrg/calculation-service`.
   - After completing these steps, go should automatically install the necessary dependencies needed to begin developing.

If set up following the Go documentation, the calculation-service folder/repository should now be found under go → src → github.com → AITestingOrg.

## Running

To run locally:
   - Have the [gmaps adapter](https://github.com/AITestingOrg/gmaps-adapter) repository cloned. Follow the instructions within the README to initiate it using the terminal. Make sure that it's listening and serving on port 8080.
   - Build any changes you may have made (or if you haven't built at all) with the command `go build` within the terminal under the calculation-service folder.
   - Begin the calculation service with the command `./calculation-service`, it should log that the service is running.
   - When running locally, Eureka will not be initialized.

To run using docker-compose: (be sure to have Docker installed)
   - Be sure to run `docker-compose pull` to make sure you have the latest images pulled.
   - Run `docker-compose up --build` within the calculation-service folder.
   - Both services gmaps and calculation should be initiated and also registered by [Eureka](https://github.com/Netflix/eureka).

## Endpoints

#### POST to api/v1/cost

With body:

```json
{
    "origin": "stringOrigin",
    "destination": "stringDestination"
    "userId": "stringUUID"
}
```
