package internal

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type table struct {
	name string
}

func newTable(name string) *table {
	table := &table{
		name: name,
	}

	return table
}

func (t *table) tableResourceBlockPath() []string {
	return []string{"aws_dynamodb_table", t.name}
}

func (t *table) resourceIdentifier() string {
	return "[aws_dynamodb_table." + t.name + ".arn]"
}

func (t *table) build(root *hclwrite.Body) {
	table := root.AppendNewBlock("resource", t.tableResourceBlockPath())
	table.Body().SetAttributeValue("name", cty.StringVal(t.name))
	table.Body().SetAttributeValue("billing_mode", cty.StringVal("PROVISIONED"))
	table.Body().SetAttributeValue("read_capacity", cty.NumberIntVal(5))
	table.Body().SetAttributeValue("write_capacity", cty.NumberIntVal(5))
	table.Body().SetAttributeValue("hash_key", cty.StringVal("id"))
	table.Body().AppendNewline()

	attribute := table.Body().AppendNewBlock("attribute", nil)
	attribute.Body().SetAttributeValue("name", cty.StringVal("id"))
	attribute.Body().SetAttributeValue("type", cty.StringVal("S"))

	root.AppendNewline()
}

type tableFunctionPermission struct {
	table       *table
	function    *function
	permissions []string
}

func newTableFunctionPermission(table *table, function *function, permissions []string) *tableFunctionPermission {
	return &tableFunctionPermission{
		table:       table,
		function:    function,
		permissions: permissions,
	}
}

func (t *tableFunctionPermission) policyDocumentName() string {
	return t.table.name + "_" + t.function.lambdaName() + "_policy_document"
}

func (t *tableFunctionPermission) policyDocumentPath() []string {
	return []string{"aws_iam_policy_document", t.policyDocumentName()}
}

func (f *tableFunctionPermission) principalIdentifier() string {
	return "[aws_iam_role." + f.function.lambdaName() + "_role.arn]"
}

func (f *tableFunctionPermission) policyName() string {
	return f.table.name + "_" + f.function.lambdaName() + "_policy"
}

func (f *tableFunctionPermission) policyResourceBlockPath() []string {
	return []string{"aws_dynamodb_resource_policy", f.policyName()}
}

func (t *tableFunctionPermission) Build(root *hclwrite.Body) {
	policyDocumentId := t.policyDocumentPath()
	blockPermission := root.AppendNewBlock("data", policyDocumentId)
	statement := blockPermission.Body().AppendNewBlock("statement", nil)
	statement.Body().SetAttributeValue("effect", cty.StringVal("Allow"))

	l := []cty.Value{}
	for _, p := range t.permissions {
		if p == "get" {
			l = append(l, cty.StringVal("dynamodb:GetItem"))
		}
		if p == "put" {
			l = append(l, cty.StringVal("dynamodb:PutItem"))
		}
	}
	statement.Body().SetAttributeValue("actions", cty.ListVal(l))
	statement.Body().AppendNewline()

	principals := statement.Body().AppendNewBlock("principals", nil)
	principals.Body().SetAttributeValue("type", cty.StringVal("AWS"))
	principals.Body().SetAttributeRaw("identifiers", hclwrite.TokensForIdentifier(t.principalIdentifier()))
	statement.Body().AppendNewline()

	statement.Body().SetAttributeRaw("resources", hclwrite.TokensForIdentifier(t.table.resourceIdentifier()))
	root.AppendNewline()

	policy := root.AppendNewBlock("resource", t.policyResourceBlockPath())
	policy.Body().SetAttributeTraversal("resource_arn", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_dynamodb_table",
		},
		hcl.TraverseAttr{
			Name: t.table.name,
		},
		hcl.TraverseAttr{
			Name: "arn",
		},
	})

	policy.Body().SetAttributeTraversal("policy", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "data",
		},
		hcl.TraverseAttr{
			Name: "aws_iam_policy_document",
		},
		hcl.TraverseAttr{
			Name: t.policyDocumentName(),
		},
		hcl.TraverseAttr{
			Name: "json",
		},
	})
	root.AppendNewline()
}
