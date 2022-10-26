package opensearchconfig_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/shopsmart/opensearchconfig"
)

var _ = Describe("Opensearchconfig", func() {

	It("Should pull the addresses environment variable", func() {
		os.Setenv("OPENSEARCH_ADDRESSES", "https://opensearch.community.dev,http://localhost:9200")

		cfg, err := opensearchconfig.ConfigFromEnv()

		Expect(err).Should(BeNil())
		Expect(cfg.Addresses).Should(ContainElements(
			"https://opensearch.community.dev",
			"http://localhost:9200",
		))
	})

})
