package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

// table.products [GET PUT]
// api.POST products [201 409 500]

func HandleRequest(ctx context.Context) error {
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
