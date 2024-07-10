package internal

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestProviderCreation(t *testing.T) {
	hclFile := hclwrite.NewEmptyFile()

	provider(hclFile.Body(), "project")

	cupaloy.SnapshotT(t, hclFile.Bytes())
}
