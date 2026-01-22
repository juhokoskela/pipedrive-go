package main

import (
	"context"
	"log"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
	v2 "github.com/juhokoskela/pipedrive-go/pipedrive/v2"
)

func main() {
	const token = "YOUR_API_TOKEN"

	client, err := v2.NewClient(pipedrive.Config{
		Auth: pipedrive.APITokenAuth(token),
	})
	if err != nil {
		log.Fatal(err)
	}

	deals, _, err := client.Deals.List(
		context.Background(),
		v2.WithDealsPageSize(10),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("deals=%d", len(deals))
}
