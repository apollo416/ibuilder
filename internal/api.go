package internal

import (
	"github.com/apollo416/ibuilder/infra"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type parentIdTraversable interface {
	traversableId() hcl.Traversal
}

type api struct {
	infraApi *infra.Api
}

func newApi(a *infra.Api) *api {
	return &api{
		infraApi: a,
	}
}

func (a *api) traversableId() hcl.Traversal {
	return hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_api_gateway_rest_api",
		},
		hcl.TraverseAttr{
			Name: a.infraApi.Name,
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	}
}

func (a *api) build(root *hclwrite.Body) {
	ap := root.AppendNewBlock("resource", []string{"aws_api_gateway_rest_api", a.infraApi.Name})
	ap.Body().SetAttributeValue("name", cty.StringVal(a.infraApi.Name))
	ap.Body().AppendNewline()
	endpointConfig := ap.Body().AppendNewBlock("endpoint_configuration", nil)
	endpointConfig.Body().SetAttributeValue("types", cty.ListVal([]cty.Value{cty.StringVal("REGIONAL")}))
	ap.Body().AppendNewline()

	lifecycle := ap.Body().AppendNewBlock("lifecycle", nil)
	lifecycle.Body().SetAttributeValue("create_before_destroy", cty.True)
	root.AppendNewline()
}

type resource struct {
	api      *api
	resource *infra.Resource
	parent   parentIdTraversable
}

func newResource(res *infra.Resource, api *api, parent parentIdTraversable) *resource {
	return &resource{
		api:      api,
		resource: res,
		parent:   parent,
	}
}

func (a *resource) traversableId() hcl.Traversal {
	return hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_api_gateway_resource",
		},
		hcl.TraverseAttr{
			Name: a.resource.Name,
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	}
}

func (r *resource) build(root *hclwrite.Body) {
	resource := root.AppendNewBlock("resource", []string{"aws_api_gateway_resource", r.resource.Name})
	resource.Body().SetAttributeTraversal("rest_api_id", r.api.traversableId())
	resource.Body().SetAttributeTraversal("parent_id", r.parent.traversableId())
	resource.Body().SetAttributeValue("path_part", cty.StringVal(r.resource.Path))
	root.AppendNewline()
}

type method struct {
	resource *resource
	function *function
	method   *infra.Method
}

func newMethod(resource *resource, function *function, m *infra.Method) *method {
	return &method{
		resource: resource,
		function: function,
		method:   m,
	}
}

func (m *method) build(root *hclwrite.Body) {
	method := root.AppendNewBlock("resource", []string{"aws_api_gateway_method", "crops_post"})
	method.Body().SetAttributeTraversal("rest_api_id", m.resource.api.traversableId())
	method.Body().SetAttributeTraversal("resource_id", m.resource.traversableId())
	method.Body().SetAttributeValue("http_method", cty.StringVal("POST"))
	method.Body().SetAttributeValue("authorization", cty.StringVal("NONE"))
	root.AppendNewline()

	integration := root.AppendNewBlock("resource", []string{"aws_api_gateway_integration", "crops_post"})
	integration.Body().SetAttributeTraversal("rest_api_id", m.resource.api.traversableId())
	integration.Body().SetAttributeTraversal("resource_id", m.resource.traversableId())
	integration.Body().SetAttributeTraversal("http_method", hcl.Traversal{hcl.TraverseRoot{
		Name: "aws_api_gateway_method"}, hcl.TraverseAttr{Name: "crops_post"}, hcl.TraverseAttr{Name: "http_method"}})
	integration.Body().SetAttributeValue("integration_http_method", cty.StringVal("POST"))
	integration.Body().SetAttributeValue("type", cty.StringVal("AWS_PROXY"))
	integration.Body().SetAttributeTraversal("uri", m.function.traversalInvokeArn())
	root.AppendNewline()

	permission := root.AppendNewBlock("resource", []string{"aws_lambda_permission", "crops_post"})
	permission.Body().SetAttributeValue("statement_id", cty.StringVal("AllowExecution"+"_crops_post_"+"FromAPI"))
	permission.Body().SetAttributeValue("action", cty.StringVal("lambda:InvokeFunction"))
	permission.Body().SetAttributeTraversal("function_name", hcl.Traversal{hcl.TraverseRoot{
		Name: "aws_lambda_function"}, hcl.TraverseAttr{Name: "crops_crop_add"}, hcl.TraverseAttr{Name: "function_name"}})
	permission.Body().SetAttributeValue("principal", cty.StringVal("apigateway.amazonaws.com"))
	lambdaIdentifier := "[aws_lambda_function.crops_crop_add]"
	permission.Body().SetAttributeRaw("depends_on", hclwrite.TokensForIdentifier(lambdaIdentifier))
	root.AppendNewline()
}
