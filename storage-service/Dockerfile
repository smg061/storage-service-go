FROM golang:1.18-alpine
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY secrets.json .
RUN go mod download
COPY . .

RUN go build -o /build

CMD [ "/build" ]