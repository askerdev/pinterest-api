FROM golang:1.22.1-alpine as builder

WORKDIR /build

COPY . .

RUN go build -o /main cmd/main.go

FROM scratch

COPY --from=builder /main /bin/main

ENV APP_MODE="PRODUCTION"
ENV APP_PORT=8081
ENV POSTGRES_URL="postgres://root:root@localhost:5432/root?sslmode=disable"

EXPOSE 8081

ENTRYPOINT ["/bin/main"]