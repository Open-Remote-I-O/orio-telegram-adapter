FROM docker.io/library/golang:1.21.6-bullseye as builder

WORKDIR /go/src

RUN mkdir build
COPY . .

# Compiling stage
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/orio-telegram-adapter ./src/cmd/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /cmd
# Retrieving the compiled binary from the stage before
COPY --from=builder /go/src/build .
CMD ["./orio-telegram-adapter"]
