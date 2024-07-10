package internal

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestVariablesCreation(t *testing.T) {
	hclFile := hclwrite.NewEmptyFile()

	variables(hclFile.Body())

	cupaloy.SnapshotT(t, hclFile.Bytes())
}
