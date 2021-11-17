FROM golang:1.16-buster as build

WORKDIR /go/src/app
ADD ./xds /go/src/app

RUN go get -d -v ./...
RUN go build \
    -o /go/bin/app \
    ./cmd/controller.go

FROM discolix/base:debug

COPY --from=build /go/bin/app /app/controller

ENTRYPOINT ["/app/controller"]