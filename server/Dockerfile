FROM golang:1.18-alpine

WORKDIR /server

COPY . .
RUN go mod download

RUN go build -o /api

EXPOSE ${PORT}

CMD [ "/api" ]
