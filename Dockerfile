FROM golang

RUN mkdir -p /go/src/calculation-service  
WORKDIR /go/src/calculation-service

ADD . /go/src/calculation-service/

RUN go get -v

EXPOSE 8000