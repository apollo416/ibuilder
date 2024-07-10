package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

// table.products [PUT]
// api.PUT products/{id} [200 500]

func HandleRequest(ctx context.Context) error {
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
