package internal

import (
	"testing"

	"github.com/apollo416/ibuilder/infra"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestApiCreation(t *testing.T) {
	hclFile := hclwrite.NewEmptyFile()
	testAPI.build(hclFile.Body())
	cupaloy.SnapshotT(t, hclFile.Bytes())
}

func TestResourceCreation(t *testing.T) {
	hclFile := hclwrite.NewEmptyFile()

	res := &infra.Resource{
		Name:    "crops",
		Path:    "crops",
		Parent:  nil,
		Methods: []*infra.Method{},
	}

	resource := newResource(res, testAPI, testAPI)
	resource.build(hclFile.Body())

	res2 := &infra.Resource{
		Name:    "crop",
		Path:    "{id}",
		Parent:  nil,
		Methods: []*infra.Method{},
	}

	resource2 := newResource(res2, testAPI, resource)
	resource2.build(hclFile.Body())

	cupaloy.SnapshotT(t, hclFile.Bytes())
}

func TestMethodCreation(t *testing.T) {
	hclFile := hclwrite.NewEmptyFile()

	infraMethod := &infra.Method{
		Function:  &infra.Function{},
		Method:    "POST",
		Responses: []string{"200"},
	}

	method := newMethod(testResource, testFunction, infraMethod)
	method.build(hclFile.Body())
	cupaloy.SnapshotT(t, hclFile.Bytes())
}
