package main

import (
	"context"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"

	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/shopsmart/opensearchconfig"
)

func main() {
	log.SetLevel(log.DebugLevel)

	client, err := opensearchconfig.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = PingOpenSearch(client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the opensearch cluster")
}

// PingOpenSearch will make a ping request to the opensearch cluster
func PingOpenSearch(client *opensearch.Client) error {
	ping := opensearchapi.PingRequest{}

	log.Debug("making a ping request to opensearch")

	resp, err := ping.Do(context.Background(), client)
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
