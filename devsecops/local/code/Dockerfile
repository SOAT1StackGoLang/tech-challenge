# build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
ADD ./code .
RUN ls -alth
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v cmd/web/*.go
RUN go build -o /go/bin/migs -v cmd/migrations/*.go

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
COPY --from=builder /go/bin/migs /migs
CMD /migs; /app
EXPOSE 8000