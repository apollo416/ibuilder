package internal

import (
	"github.com/apollo416/ibuilder/infra"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

const (
	functionHandlerName = "bootstrap"
)

type function struct {
	function *infra.Function
	hash     string
	zipPath  string
}

func newFunction(f *infra.Function, hash, zipPath string) *function {
	return &function{
		function: f,
		hash:     hash,
		zipPath:  zipPath,
	}
}

func (f *function) lambdaName() string {
	return f.function.Name + "_function"
}

func (f *function) lambdaRoleName() string {
	return f.lambdaName() + "_role"
}

func (f *function) traversalInvokeArn() hcl.Traversal {
	return hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_lambda_function",
		},
		hcl.TraverseAttr{
			Name: f.lambdaName(),
		},
		hcl.TraverseAttr{
			Name: "invoke_arn",
		},
	}
}

func (f *function) build(root *hclwrite.Body) {
	function := root.AppendNewBlock("resource", []string{"aws_lambda_function", f.lambdaName()})
	function.Body().SetAttributeValue("filename", cty.StringVal(f.zipPath))
	function.Body().SetAttributeValue("function_name", cty.StringVal(f.function.Name))
	function.Body().SetAttributeValue("runtime", cty.StringVal("provided.al2023"))
	function.Body().SetAttributeValue("handler", cty.StringVal(functionHandlerName))
	function.Body().SetAttributeValue("timeout", cty.NumberIntVal(10))
	function.Body().SetAttributeValue("memory_size", cty.NumberIntVal(128))
	function.Body().SetAttributeValue("publish", cty.True)
	function.Body().SetAttributeValue("reserved_concurrent_executions", cty.NumberIntVal(-1))
	function.Body().SetAttributeValue("architectures", cty.ListVal([]cty.Value{cty.StringVal("arm64")}))
	function.Body().SetAttributeValue("source_code_hash", cty.StringVal(f.hash))

	function.Body().SetAttributeTraversal("role", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_iam_role",
		},
		hcl.TraverseAttr{
			Name: f.lambdaRoleName(),
		},
		hcl.TraverseAttr{
			Name: "arn",
		},
	})
	root.AppendNewline()

	role := root.AppendNewBlock("resource", []string{"aws_iam_role", f.lambdaRoleName()})
	role.Body().SetAttributeValue("name", cty.StringVal(f.lambdaRoleName()))
	role.Body().SetAttributeTraversal("assume_role_policy", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "data",
		},
		hcl.TraverseAttr{
			Name: "aws_iam_policy_document",
		},
		hcl.TraverseAttr{
			Name: f.lambdaRoleName(),
		},
		hcl.TraverseAttr{
			Name: "json",
		},
	})
	root.AppendNewline()

	data := root.AppendNewBlock("data", []string{"aws_iam_policy_document", f.lambdaRoleName()})

	statement := data.Body().AppendNewBlock("statement", nil)
	statement.Body().SetAttributeValue("actions", cty.ListVal([]cty.Value{cty.StringVal("sts:AssumeRole")}))

	principals := statement.Body().AppendNewBlock("principals", nil)
	principals.Body().SetAttributeValue("type", cty.StringVal("Service"))
	principals.Body().SetAttributeValue("identifiers", cty.ListVal([]cty.Value{cty.StringVal("lambda.amazonaws.com")}))

	root.AppendNewline()
}
