package main

import (
	"context"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"

	opensearch "github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/shopsmart/opensearchconfig-go"
)

func main() {
	log.SetLevel(log.DebugLevel)

	ctx := context.Background()

	client, err := opensearchconfig.NewClientFromEnv(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = PingOpenSearch(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the opensearch cluster")
}

// PingOpenSearch will make a ping request to the opensearch cluster
func PingOpenSearch(ctx context.Context, client *opensearch.Client) error {
	ping := opensearchapi.PingReq{}

	log.Debug("making a ping request to opensearch")

	req, err := ping.GetRequest()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		log.WithFields(log.Fields{"status": resp.Status()}).Debug("ping response status")

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unable to read response body from ping request")
		}

		err = fmt.Errorf("failed to ping opensearch cluster")

		log.WithFields(log.Fields{"body": respBody}).Debug("ping response body")
		return err
	}

	return nil
}
