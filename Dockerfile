FROM golang

LABEL maintainer="erinswyoo@gmail.com, lamschan996@ufl.edu"
RUN mkdir -p /go/src/calculation-service  
WORKDIR /go/src/calculation-service

ADD . /go/src/calculation-service/

RUN go get -v

EXPOSE 8000