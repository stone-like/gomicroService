package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

var ErrParseCert = errors.New("failed to parse root certificate")

type TLSConfig struct {
	CertFile      string
	KeyFile       string
	CAFile        string
	ServerAddress string
	Server        bool
}

func SetupTLSConfig(cfg TLSConfig) (*tls.Config, error) {

	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS13}
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig.Certificates[0] = cert
	}

	if cfg.CAFile != "" {
		b, err := os.ReadFile(cfg.CAFile)
		if err != nil {
			return nil, err
		}
		ca := x509.NewCertPool()
		if ok := ca.AppendCertsFromPEM([]byte(b)); !ok {
			return nil, ErrParseCert
		}

		tlsConfig.ServerName = cfg.ServerAddress

		if cfg.Server {
			//サーバーのConfigではClientAuthを設定し、サーバーがクライアントの証明書を検証、またクライアントがサーバの証明書を検証できるようにしている(mTLS)
			tlsConfig.ClientCAs = ca
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		} else {
			//クライアントだったらRootCAを設定することでmTLSをできるようにしている
			tlsConfig.RootCAs = ca
		}
	}

	return tlsConfig, nil
}
