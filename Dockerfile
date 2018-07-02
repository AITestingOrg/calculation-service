FROM golang

LABEL maintainer="erinswyoo@gmail.com, lamschan996@ufl.edu"
RUN go env GOPATH
RUN mkdir -p /go/src/github.com/AITestingOrg/calculation-service
WORKDIR /go/src/github.com/AITestingOrg/calculation-service

COPY . .

RUN go get -v ./...
RUN go build

EXPOSE 8080

CMD ["calculation-service"]