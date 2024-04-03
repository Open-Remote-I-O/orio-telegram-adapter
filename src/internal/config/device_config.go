package config

const (
	certPathDefault       = "/etc/ssl/certs/orio-server.crt"
	privateKeyPathDefault = "/etc/ssl/private/orio-server.key"
	caCertPathDefault     = "/etc/ssl/certs/orio-ca.crt"
)

type DeviceConfig struct {
	Orio_tls_cert_path string
	Orio_tls_key_path  string
	Orio_ca_cert_path  string
}

func NewDeviceConfig() DeviceConfig {
	var deviceConf DeviceConfig
	deviceConf.Orio_tls_cert_path = GetEnvOrDefault("ORIO_TLS_CERT_PATH", certPathDefault)

	deviceConf.Orio_tls_key_path = GetEnvOrDefault("ORIO_TLS_KEY_PATH", privateKeyPathDefault)

	deviceConf.Orio_ca_cert_path = GetEnvOrDefault("ORIO_CA_TLS_CERT_PATH", caCertPathDefault)

	return deviceConf
}
