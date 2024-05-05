FROM golang:alpine AS builder
WORKDIR $GOPATH/src/app/
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app cmd/example/*.go

FROM alpine
COPY --from=builder /go/bin/app /go/bin/app

ENTRYPOINT [ "/go/bin/app" ]
