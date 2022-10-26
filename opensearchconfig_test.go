package opensearchconfig_test

import (
	"context"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/shopsmart/opensearchconfig"
)

var _ = Describe("Opensearchconfig", func() {

	var (
		ctx = context.Background()
	)

	BeforeEach(func() {
		// Unset all OPENSEARCH_ environment variables
		for _, pair := range os.Environ() {
			if strings.HasPrefix(pair, "OPENSEARCH_") {
				os.Unsetenv(strings.Split(pair, "=")[0])
			}
		}
	})

	// This is built into the opensearch-go library, we are simply validating that we are not overwriting the functionality
	It("Should pull the OPENSEARCH_URL environment variable", func() {
		os.Setenv("OPENSEARCH_URL", "https://opensearch.community.dev,http://localhost:9200")

		cfg, err := opensearchconfig.ConfigFromEnv(ctx)

		Expect(err).Should(BeNil())
		Expect(cfg.Addresses).Should(BeEmpty())
	})

	It("Should pull the OPENSEARCH_SKIP_SSL environment variable", func() {
		os.Setenv("OPENSEARCH_SKIP_SSL", "true")

		cfg, err := opensearchconfig.ConfigFromEnv(ctx)

		Expect(err).Should(BeNil())
		Expect(cfg.Transport).ShouldNot(BeNil())
	})

	Context("Auth", func() {
		It("Should not configure anything with none auth", func() {
			os.Setenv("OPENSEARCH_AUTH", "none")

			cfg, err := opensearchconfig.ConfigFromEnv(ctx)

			Expect(err).Should(BeNil())
			Expect(cfg.Username).Should(BeEmpty())
			Expect(cfg.Password).Should(BeEmpty())
			Expect(cfg.Signer).Should(BeNil())
		})

		It("Should configure the username and password with basic auth", func() {
			os.Setenv("OPENSEARCH_AUTH", "basic")
			os.Setenv("OPENSEARCH_USERNAME", "username")
			os.Setenv("OPENSEARCH_PASSWORD", "password")

			cfg, err := opensearchconfig.ConfigFromEnv(ctx)

			Expect(err).Should(BeNil())
			Expect(cfg.Username).Should(Equal("username"))
			Expect(cfg.Password).Should(Equal("password"))
			Expect(cfg.Signer).Should(BeNil())
		})

		It("Should error out if the username is not provided with basic auth", func() {
			os.Setenv("OPENSEARCH_AUTH", "basic")
			os.Setenv("OPENSEARCH_PASSWORD", "password")

			_, err := opensearchconfig.ConfigFromEnv(ctx)

			Expect(err).Should(Equal(opensearchconfig.ErrMissingCredentials))
		})

		It("Should error out if the password is not provided with basic auth", func() {
			os.Setenv("OPENSEARCH_AUTH", "basic")
			os.Setenv("OPENSEARCH_USERNAME", "username")

			_, err := opensearchconfig.ConfigFromEnv(ctx)

			Expect(err).Should(Equal(opensearchconfig.ErrMissingCredentials))
		})

		It("Should configure the aws signer with iam auth", func() {
			os.Setenv("OPENSEARCH_AUTH", "iam")

			cfg, err := opensearchconfig.ConfigFromEnv(ctx)

			Expect(err).Should(BeNil())
			Expect(cfg.Username).Should(BeEmpty())
			Expect(cfg.Password).Should(BeEmpty())
			Expect(cfg.Signer).ShouldNot(BeNil())
		})
	})
})
