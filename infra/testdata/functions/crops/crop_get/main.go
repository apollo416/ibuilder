package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

// table.crops [GET]
// api.GET crops/{id} [200 404 500]

func HandleRequest(ctx context.Context) error {
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
