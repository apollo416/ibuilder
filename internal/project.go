package internal

import (
	"github.com/apollo416/ibuilder/infra"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type project struct {
	api     *api
	project *infra.Project
}

func newProject(a *api, p *infra.Project) *project {
	project := &project{
		api:     a,
		project: p,
	}

	return project
}

func (p *project) build(root *hclwrite.Body) {

	variables(root)
	terraform(root, "buckets3", "dynamotable")
	provider(root, p.project.Name)

	for _, x := range p.project.Functions {
		p := newFunction(x, "", "")
		p.build(root)
	}

	for _, x := range p.project.Tables {
		t := newTable(x.Name)
		t.build(root)
	}

	for _, x := range p.project.Apis {
		a := newApi(x)
		a.build(root)

		for _, w := range x.Resources {
			r := newResource(w, p.api, p.api)
			r.build(root)
			for _, y := range w.Methods {
				p := newFunction(y.Function, "", "")
				m := newMethod(r, p, y)
				m.build(root)
			}
		}
	}
}
