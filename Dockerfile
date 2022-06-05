FROM golang:1.18 as build-env

WORKDIR /go/src/app

ADD . ./

RUN go build -ldflags="-s -w" -o /go/bin/famed-backend

FROM gcr.io/distroless/base

COPY --from=build-env /go/bin/famed-backend /go/src/app/config.json /

CMD ["/famed-backend"]
