package internal

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestFunctionCreation(t *testing.T) {
	hclFile := hclwrite.NewEmptyFile()

	testFunction.build(hclFile.Body())

	cupaloy.SnapshotT(t, hclFile.Bytes())
}
