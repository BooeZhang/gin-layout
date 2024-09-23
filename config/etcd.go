package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// Etcd etcd配置项
type Etcd struct {
	Endpoints            []string `json:"endpoints"               mapstructure:"endpoints"`
	Timeout              int      `json:"timeout"                 mapstructure:"timeout"`
	RequestTimeout       int      `json:"request-timeout"         mapstructure:"request-timeout"`
	LeaseExpire          int      `json:"lease-expire"            mapstructure:"lease-expire"`
	Username             string   `json:"username"                mapstructure:"username"`
	Password             string   `json:"password"                mapstructure:"password"`
	UseTLS               bool     `json:"use-tls"                 mapstructure:"use-tls"`
	CaCert               string   `json:"ca-cert"                 mapstructure:"ca-cert"`
	Cert                 string   `json:"cert"                    mapstructure:"cert"`
	Key                  string   `json:"key"                     mapstructure:"key"`
	HealthBeatPathPrefix string   `json:"health_beat_path_prefix" mapstructure:"health_beat_path_prefix"`
	HealthBeatIFaceName  string   `json:"health_beat_iface_name"  mapstructure:"health_beat_iface_name"`
	Namespace            string   `json:"namespace"               mapstructure:"namespace"`
}

// Validate verifies flags passed to Redis.
func (o *Etcd) Validate() []error {
	errs := []error{}

	if len(o.Endpoints) == 0 {
		errs = append(errs, fmt.Errorf("etcd endpoints can not be empty"))
	}

	if o.RequestTimeout <= 0 {
		errs = append(errs, fmt.Errorf("--etcd.request-timeout cannot be negative"))
	}

	return errs
}

// GetEtcdTLSConfig returns etcd tls config.
func (o *Etcd) GetEtcdTLSConfig() (*tls.Config, error) {
	var (
		cert       tls.Certificate
		certLoaded bool
		capool     *x509.CertPool
	)
	if o.Cert != "" && o.Key != "" {
		var err error
		cert, err = tls.LoadX509KeyPair(o.Cert, o.Key)
		if err != nil {
			return nil, err
		}
		certLoaded = true
		o.UseTLS = true
	}
	if o.CaCert != "" {
		data, err := ioutil.ReadFile(o.CaCert)
		if err != nil {
			return nil, err
		}
		capool = x509.NewCertPool()
		for {
			var block *pem.Block
			block, _ = pem.Decode(data)
			if block == nil {
				break
			}
			cacert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			capool.AddCert(cacert)
		}
		o.UseTLS = true
	}

	if o.UseTLS {
		// nolint: gosec
		cfg := &tls.Config{
			RootCAs:            capool,
			InsecureSkipVerify: false,
		}
		if certLoaded {
			cfg.Certificates = []tls.Certificate{cert}
		}

		return cfg, nil
	}

	return nil, nil
}
