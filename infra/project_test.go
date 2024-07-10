package infra

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestProjectFromFolder(t *testing.T) {
	project := ProjectFromFolder(TestConfig.BaseDir)

	if project.Name != "testdata" {
		t.Errorf("project.Name = %s; want testdata", project.Name)
	}

	v, _ := json.MarshalIndent(project, "", "  ")
	fmt.Println(string(v))
}
