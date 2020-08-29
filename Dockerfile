FROM golang

COPY . /go/src/app

WORKDIR /go/src/app



RUN go get -v ./api
RUN go install ./api


