FROM docker.io/library/golang:1.24.0-bullseye as builder

ARG ORIO_SERVER_CERT_PATH
ARG ORIO_SERVER_KEY_PATH
ARG ORIO_CA_CERT_PATH

RUN go install github.com/Open-Remote-I-O/cert_gen_cli@v0.3.1

WORKDIR /go/src

RUN mkdir build

COPY go.mod . 
COPY go.sum .

RUN go mod download

COPY . .

# Compiling stage
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./build/orio-telegram-adapter ./cmd/main.go

FROM gcr.io/distroless/static-debian12

# find better approach
ARG ORIO_SERVER_CERT_PATH
ARG ORIO_SERVER_KEY_PATH
ARG ORIO_CA_CERT_PATH

WORKDIR /cmd
# Retrieving the compiled binary from the stage before
COPY --from=builder /go/src/build/orio-telegram-adapter ./orio-telegram-adaptersrc

COPY ${ORIO_SERVER_CERT_PATH} /etc/ssl/certs/orio-server.crt
COPY ${ORIO_SERVER_KEY_PATH} /etc/ssl/private/orio-server.key
COPY ${ORIO_CA_CERT_PATH} /etc/ssl/certs/orio-ca.crt

CMD ["./orio-telegram-adapter"]
