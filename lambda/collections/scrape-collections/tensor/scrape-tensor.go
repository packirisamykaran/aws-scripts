package tensor

import (
	"context"
	"fmt"
	"log"

	"github.com/apollographql/apollo"
	"github.com/apollographql/operation"
)

const API_KEY = "eff00cb4-025d-4bf5-84d5-eb170054af72"

func getMintList(slug string) {
	// Create an Apollo Client instance with the GraphQL endpoint URL
	client := apollo.NewClient("https://api.tensor.so/graphql", nil)

	// Construct the operation
	op := &operation.Operation{
		Query: `
			query MintList($slug: String!, $limit: Int, $after: String) {
				mintList(slug: $slug, limit: $limit, after: $after)
			}
		`,
		Variables: map[string]interface{}{
			"slug":  slug,
			"limit": nil, // Optional: Set your desired limit value
			"after": nil, // Optional: Set your desired after value
		},
	}

	// Create a context and execute the operation
	ctx := context.Background()
	var resp *apollo.Response
	err := client.Execute(ctx, op, &resp)
	if err != nil {
		log.Fatal("GraphQL query failed:", err)
	}

	// Access the results
	mintList := resp.Data.Get("mintList").String()
	fmt.Println(mintList)
}

func ScrapeTensor() {
	getMintList("madlads")
}
