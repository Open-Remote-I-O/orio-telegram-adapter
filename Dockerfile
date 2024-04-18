FROM docker.io/library/golang:1.22.2-bullseye as builder

ARG orioServerCertPath
ARG orioServerKeyPath
ARG orioCaCertPath

RUN go install github.com/Open-Remote-I-O/cert_gen_cli@32568dc

WORKDIR /go/src

RUN mkdir build

COPY go.mod . 
COPY go.sum .

RUN go mod download

COPY . .

# Compiling stage
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./build/orio-telegram-adapter ./src/cmd/main.go

FROM gcr.io/distroless/static-debian12

ARG orioServerCertPath
ARG orioServerKeyPath
ARG orioCaCertPath

WORKDIR /cmd
# Retrieving the compiled binary from the stage before
COPY --from=builder /go/src/build .

COPY ${orioServerCertPath} /etc/ssl/certs/orio-server.crt
COPY ${orioServerKeyPath} /etc/ssl/private/orio-server.key
COPY ${orioCaCertPath} /etc/ssl/certs/orio-ca.crt

CMD ["./orio-telegram-adapter"]
