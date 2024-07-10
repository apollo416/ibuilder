package internal

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestTableCreation(t *testing.T) {
	hclFile := hclwrite.NewEmptyFile()
	testTable.build(hclFile.Body())
	cupaloy.SnapshotT(t, hclFile.Bytes())
}

func TestTableFunctionPermissionCreation(t *testing.T) {
	hclFile := hclwrite.NewEmptyFile()

	tfp := newTableFunctionPermission(
		testTable,
		testFunction,
		[]string{"get", "put"},
	)

	tfp.Build(hclFile.Body())
	cupaloy.SnapshotT(t, hclFile.Bytes())
}
