FROM docker.io/library/golang:1.22.2-alpine3.19

RUN go install github.com/Open-Remote-I-O/cert_gen_cli@v0.3.0

RUN go install github.com/cosmtrek/air@v1.51.0

WORKDIR /app

COPY go.mod . 
COPY go.sum .

COPY local-cert-gen-conf.toml .

RUN go mod download 

# Add you're CA authority certificate in order to allow mTLS communication
COPY ca.crt /etc/ssl/certs/orio-ca.crt
COPY ca.key /etc/ssl/private/orio-ca.key

RUN cert_gen_cli genCaParentCert \ 
	--key-out /etc/ssl/private/ \ 
	--cert-out /etc/ssl/certs/ \ 
	-n mock-device \ 
	-c /etc/ssl/certs/orio-ca.crt \ 
	-k /etc/ssl/private/orio-ca.key \
	--config-file-path ./local-cert-gen-conf.toml

CMD ["air", "-c", "./mock/device/device-mock.air.toml"]
