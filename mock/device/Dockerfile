FROM docker.io/library/golang:1.22.0-bullseye

RUN go install github.com/Open-Remote-I-O/cert_gen_cli@v0.1.4

RUN go install github.com/cosmtrek/air@v1.49.0

WORKDIR /app

# Add you're CA authority certificate in order to allow mTLS communication
COPY ca.crt /etc/ssl/certs/orio-ca.crt
COPY ca.key /etc/ssl/private/orio-ca.key

RUN cert_gen_cli genCaParentCert --organization-name orio-local -o /etc/ssl/certs/ -n mock-device -c /etc/ssl/certs/orio-ca.crt -k /etc/ssl/private/orio-ca.key

RUN mv /etc/ssl/certs/mock-device.key /etc/ssl/private/mock-device.key

CMD ["air", "-c", "./mock/device/device-mock.air.toml"]
