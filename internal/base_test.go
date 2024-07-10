package internal

import "github.com/apollo416/ibuilder/infra"

var (
	testAPI      *api
	testTable    *table
	testResource *resource
	testProject  *project
	testFunction *function
)

func init() {

	infraApi := &infra.Api{
		Name:      "crops",
		Resources: []*infra.Resource{},
	}

	cropsTable := &infra.Table{
		Name: "crops",
		Permissions: []*infra.TableFunctionPermision{
			&infra.TableFunctionPermision{},
		},
	}

	cropAddFunction := &infra.Function{
		Name:       "crop_add",
		SourcePath: "path",
	}

	res := &infra.Resource{
		Name:   "crops",
		Path:   "crops",
		Parent: nil,
		Methods: []*infra.Method{
			{
				Function:  cropAddFunction,
				Method:    "POST",
				Responses: []string{"200", "500"},
			},
		},
	}

	infraApi.Resources = []*infra.Resource{res}

	p := &infra.Project{
		Name:      "farm",
		Apis:      []*infra.Api{infraApi},
		Tables:    []*infra.Table{cropsTable},
		Functions: []*infra.Function{cropAddFunction},
		Resources: []*infra.Resource{res},
	}

	// ----------------------------------------------------

	testTable = newTable(cropsTable.Name)

	testFunction = newFunction(
		cropAddFunction,
		"hash",
		"path",
	)

	testAPI = newApi(infraApi)

	testResource = newResource(res, testAPI, testAPI)

	testProject = newProject(testAPI, p)
}
