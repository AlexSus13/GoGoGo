FROM golang:1.16.13-alpine3.14

RUN mkdir /appserv

WORKDIR /appserv

COPY . /appserv

RUN go mod download

WORKDIR /appserv/cmd

RUN go build -o binfileservergo

CMD [ "./binfileservergo" ]
