FROM golang

LABEL maintainer="erinswyoo@gmail.com, lamschan996@ufl.edu"
RUN mkdir -p /go/src/github.com/AITestingOrg/calculation-service
WORKDIR /go/src/github.com/AITestingOrg/calculation-service

ADD . /go/src/github.com/AITestingOrg/calculation-service

RUN go get -v

EXPOSE 8080