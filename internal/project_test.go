package internal

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestProjectCreation(t *testing.T) {
	hclFile := hclwrite.NewEmptyFile()
	testProject.build(hclFile.Body())
	cupaloy.SnapshotT(t, hclFile.Bytes())
}
