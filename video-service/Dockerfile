FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
# RUN go mod download
COPY ./videos ./videos
COPY *.go ./

RUN go build -o /build

CMD [ "/build" ]