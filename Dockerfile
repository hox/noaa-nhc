FROM golang:1.16-alpine

LABEL org.opencontainers.image.source https://github.com/hox/noaa-nhc

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /noaa-nhc

CMD [ "/noaa-nhc" ]