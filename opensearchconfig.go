package opensearchconfig

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go-v2/config"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

const (
	// AuthNone is the auth type to use if no auth is configured on the cluster
	AuthNone = "none"
	// AuthBasic is the auth type to use if using the internal user database on the cluster
	AuthBasic = "basic"
	// AuthIAM is the auth type to use if not using the internal user database on the cluster
	AuthIAM = "iam"
	// ENV variable to control whether traces are send to datadog
	DD_TRACE_ENABLED = "DD_TRACE_ENABLED"
	// Identify the name of the trace in datadog
	Opensearch = "opensearch"
)

var (
	// ErrMissingCredentials will be thrown when basic auth is configured but either username or password is not available
	ErrMissingCredentials = errors.New("basic auth has been set but username or password is missing")
)

// Config represents the configuration options available with this package
type Config struct {
	// Skips ssl if true
	SkipSSL bool // OPENSEARCH_SKIP_SSL
	// The auth type to use.  Options are none, basic, iam
	Auth string // OPENSEARCH_AUTH
	// The username if auth is basic
	Username string // OPENSEARCH_USERNAME
	// The password if auth is basic
	Password string // OPENSEARCH_PASSWORD
}

// GetConfig will get the Config object from the environment
func GetConfig() (Config, error) {
	v := viper.New()

	// Do not skip ssl by default
	v.SetDefault("skip-ssl", false)

	// Prefix all environment variables with OPERNSEARCH_
	v.SetEnvPrefix("opensearch")
	v.AutomaticEnv()

	// We do not want to load any files
	err := v.ReadConfig(bytes.NewBuffer([]byte("")))
	if err != nil {
		return Config{}, err
	}

	return Config{
		SkipSSL:  v.GetBool("skip_ssl"),
		Auth:     v.GetString("auth"),
		Username: v.GetString("username"),
		Password: v.GetString("password"),
	}, nil
}

// ConfigFromEnv creates an OpenSearch config object from environment variables
func ConfigFromEnv(ctx context.Context) (opensearch.Config, error) {
	opensearchConfig := opensearch.Config{}
	// Add datadog tracing
	opensearchConfig.Transport = &http.Transport{}

	cfg, err := GetConfig()
	if err != nil {
		return opensearchConfig, err
	}

	if cfg.SkipSSL {
		log.Debug("skipping ssl")
		opensearchConfig.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	if os.Getenv(DD_TRACE_ENABLED) == "true" {
		var opts []httptrace.RoundTripperOption
		opts = append(opts, httptrace.RTWithResourceNamer(func(req *http.Request) string {
			return Opensearch
		}))

		opensearchConfig.Transport = httptrace.WrapRoundTripper(
			opensearchConfig.Transport,
			opts...,
		)
	}

	switch cfg.Auth {
	case AuthBasic:
		log.Debug("configuring basic auth")

		if cfg.Username == "" || cfg.Password == "" {
			return opensearchConfig, ErrMissingCredentials
		}

		opensearchConfig.Username = cfg.Username
		opensearchConfig.Password = cfg.Password

	case AuthIAM:
		log.Debug("configuring iam auth")

		awsConfig, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return opensearchConfig, err
		}

		// Create an AWS request Signer and load AWS configuration using default config folder or env vars.
		// See https://docs.aws.amazon.com/opensearch-service/latest/developerguide/request-signing.html#request-signing-go
		signer, err := requestsigner.NewSigner(awsConfig)
		if err != nil {
			return opensearchConfig, err
		}

		opensearchConfig.Signer = signer

	case "none":
	default:
		log.Debug("disabling auth")
	}

	return opensearchConfig, nil
}

// NewClientFromEnv creates an OpenSearch client pulling configurations from the environment
func NewClientFromEnv(ctx context.Context) (*opensearch.Client, error) {
	cfg, err := ConfigFromEnv(ctx)
	if err != nil {
		return nil, err
	}

	return opensearch.NewClient(cfg)
}
