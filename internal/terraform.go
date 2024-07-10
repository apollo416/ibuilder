package internal

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func terraform(root *hclwrite.Body, s3Name, dynamodbName string) {
	terraform := root.AppendNewBlock("terraform", nil)
	terraformBody := terraform.Body()
	terraformBody.SetAttributeValue("required_version", cty.StringVal("1.9.1"))
	terraformBody.AppendNewline()

	backend := terraformBody.AppendNewBlock("backend", []string{"s3"})
	backendBody := backend.Body()
	backendBody.SetAttributeValue("bucket", cty.StringVal(s3Name))
	backendBody.SetAttributeValue("key", cty.StringVal("terraform.state"))
	backendBody.SetAttributeValue("region", cty.StringVal("us-east-1"))
	backendBody.SetAttributeValue("dynamodb_table", cty.StringVal(dynamodbName))
	backendBody.AppendNewline()

	requiredProviders := terraformBody.AppendNewBlock("required_providers", nil)
	requiredProvidersBody := requiredProviders.Body()
	requiredProvidersBody.SetAttributeValue("aws", cty.ObjectVal(map[string]cty.Value{
		"source":  cty.StringVal("hashicorp/aws"),
		"version": cty.StringVal("~> 5.0"),
	}))

	root.AppendNewline()
}
