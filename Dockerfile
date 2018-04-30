FROM golang

LABEL maintainer="erinswyoo@gmail.com, lamschan996@ufl.edu"
RUN go env GOPATH
RUN mkdir -p /go/src/calculation-service
WORKDIR /go/src/calculation-service

COPY . .

RUN go get -v ./...
RUN go build

EXPOSE 8080

CMD ["calculation-service"]