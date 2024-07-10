package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

// table.crops [PUT]
// api.PUT crops/{id} [201 500]

func HandleRequest(ctx context.Context) error {
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
