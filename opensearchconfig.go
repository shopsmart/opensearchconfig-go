package opensearchconfig

import (
	"bytes"
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go-v2/config"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
)

const (
	// DefaultAddress is the default address to use if no address is provided
	DefaultAddress = "http://localhost:9200"
)

// ConfigFromEnv creates an OpenSearch config object from environment variables
func ConfigFromEnv() (opensearch.Config, error) {
	v := viper.New()

	v.SetDefault("addresses", DefaultAddress)
	v.SetDefault("skip-ssl", false)

	v.SetEnvPrefix("opensearch")
	v.AutomaticEnv()

	err := v.ReadConfig(bytes.NewBuffer([]byte("")))
	if err != nil {
		return opensearch.Config{}, err
	}

	addresses := strings.Split(v.GetString("addresses"), ",")
	log.WithFields(log.Fields{"addresses": addresses}).Debug("configuring opensearch addresses")
	cfg := opensearch.Config{
		Addresses: addresses,
	}

	if v.GetBool("skip-ssl") {
		log.Debug("skipping ssl")
		cfg.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	switch v.GetString("auth") {
	case "iam":
		log.Debug("configuring iam auth")

		ctx := context.Background()

		awsConfig, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return cfg, err
		}

		// Create an AWS request Signer and load AWS configuration using default config folder or env vars.
		// See https://docs.aws.amazon.com/opensearch-service/latest/developerguide/request-signing.html#request-signing-go
		signer, err := requestsigner.NewSigner(awsConfig)
		if err != nil {
			return cfg, err
		}

		cfg.Signer = signer

	case "none":
	default:
		log.Debug("disabling auth")
	}

	return cfg, nil
}

// NewClientFromEnv creates an OpenSearch client pulling configurations from the environment
func NewClientFromEnv() (*opensearch.Client, error) {
	cfg, err := ConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return opensearch.NewClient(cfg)
}
