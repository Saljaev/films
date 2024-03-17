FROM golang:1.22-alpine AS builder

WORKDIR /go/src/backend

RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc git

COPY . /go/src/backend
RUN go mod download

RUN go build -o backend cmd/backend/*.go

FROM alpine

COPY --from=builder /go/src/backend/backend /backend
COPY --from=builder /go/src/backend/config/config.yaml /config/config.yaml
COPY --from=builder /go/src/backend/internal/migrations/* /migrations/

EXPOSE 8080

ENTRYPOINT ["/backend"]
CMD [""]
