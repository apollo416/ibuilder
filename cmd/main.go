package main

import (
	"fmt"
	"os"

	"github.com/apollo416/ibuilder/infra"
)

func main() {
	projectFolder := "."

	if len(os.Args) > 1 {
		projectFolder = os.Args[1]
	}

	project := infra.ProjectFromFolder(projectFolder)
	fmt.Println(project)
}
