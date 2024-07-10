package internal

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func variables(root *hclwrite.Body) {
	variable := root.AppendNewBlock("variable", []string{"workspace_iam_role"})
	variableBody := variable.Body()
	variableBody.SetAttributeRaw(
		"type",
		hclwrite.TokensForIdentifier("string"),
	)
	variableBody.SetAttributeValue("description", cty.StringVal("the iam role to assume for the workspace"))
	variableBody.SetAttributeValue("nullable", cty.False)
	variableBody.SetAttributeValue("sensitive", cty.True)
	root.AppendNewline()
}
