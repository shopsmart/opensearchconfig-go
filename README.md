# OpenSearch Config

Pulls environment variables to configure the opensearch-go library

## Usage

```go
package main

import (
	"context"

    "github.com/shopsmart/opensearchconfig"
)

func main() {
	ctx := context.Background()

	client, err := opensearchconfig.NewClientFromEnv(ctx)
	if err != nil {
        panic(err) // Beware the scary panic
	}

    // Do stuff with the client
}
```

## Configuring

### Available environment variables

- `OPENSEARCH_URL` is parsed by the opensearch-go library and configures the various addresses to be used with the client.  Separate endpoints with a comma.
- `OPENSEARCH_SKIP_SSL` set to true will enable skipping the SSL check
- `OPENSEARCH_AUTH` sets what type of auth to use for connecting to the OpenSearch cluster.  By default, it will use `none` auth.  Options are: `none`, `basic`, `iam`.
- `OPENSEARCH_USERNAME` represents the username for the user to login with only if auth is set to `basic`
- `OPENSEARCH_PASSWORD` represents the password for the user to login with only if auth is set to `basic`
