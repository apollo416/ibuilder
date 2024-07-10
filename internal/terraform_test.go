package internal

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestTerraformCreation(t *testing.T) {
	hclFile := hclwrite.NewEmptyFile()

	terraform(hclFile.Body(), "s3Name", "dynamodbName")

	cupaloy.SnapshotT(t, hclFile.Bytes())

}
