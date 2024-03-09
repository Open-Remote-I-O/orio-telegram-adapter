FROM docker.io/library/golang:1.22.1-bullseye as builder

RUN go install github.com/Open-Remote-I-O/cert_gen_cli@v0.1.4

# Add you're CA authority certificate in order to allow mTLS communication
COPY ca.crt /etc/ssl/certs/orio-ca.crt
COPY ca.key /etc/ssl/private/orio-ca.key

RUN cert_gen_cli genCaParentCert --organization-name orio -o /etc/ssl/certs/ -n orio-server -c /etc/ssl/certs/orio-ca.crt -k /etc/ssl/private/orio-ca.key

WORKDIR /go/src

RUN mkdir build

COPY . .

RUN go mod download

# Compiling stage
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/orio-telegram-adapter ./src/cmd/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /cmd
# Retrieving the compiled binary from the stage before
COPY --from=builder /go/src/build .
COPY --from=builder /etc/ssl/certs/orio-server.crt /etc/ssl/certs/orio-server.crt
COPY --from=builder /etc/ssl/certs/orio-server.key /etc/ssl/certs/orio-server.key
CMD ["./orio-telegram-adapter"]
