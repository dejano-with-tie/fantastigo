FROM golang:1.20 AS golang-plugins

WORKDIR /go/src/
COPY ./plugins/ ./
RUN go build -o /telegraf-plugins/gps gps/gps.go

FROM telegraf:1.24.4

COPY --from=golang-plugins /telegraf-plugins/gps /usr/local/bin
RUN echo -e "$(ls /usr/local/bin)"