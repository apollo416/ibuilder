package infra

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Project struct {
	Name      string      `json:"name"`
	Apis      []*Api      `json:"apis"`
	Tables    []*Table    `json:"tables"`
	Functions []*Function `json:"functions"`
	Resources []*Resource `json:"resources"`
}

func ProjectFromFolder(basePath string) *Project {
	config := NewConfig(basePath)
	projectBuilder := newProjectBuilder(config)
	return projectBuilder.build()
}

type projectBuilder struct {
	config       Config
	project      *Project
	apiFunctions map[string][]*Function
}

func newProjectBuilder(config Config) *projectBuilder {
	return &projectBuilder{
		config:       config,
		project:      &Project{},
		apiFunctions: map[string][]*Function{},
	}
}

func (b *projectBuilder) build() *Project {
	projectName := filepath.Base(b.config.BaseDir)

	b.project.Name = projectName

	b.project.Tables = b.loadTables()
	b.project.Functions = b.loadFunctions()

	b.project.Apis = b.loadApis()

	b.setTableFunctionPermissions()

	return b.project
}

func (b *projectBuilder) loadApis() []*Api {
	apis := []*Api{}

	apiFolders := listFolders(b.config.FunctionsSourceDir)

	for _, apiFolder := range apiFolders {
		api := &Api{
			Name:      apiFolder,
			Resources: []*Resource{},
		}

		api.Resources = b.loadResources(apiFolder)

		apis = append(apis, api)
	}

	return apis
}

func (b *projectBuilder) setTableFunctionPermissions() {
	reg := regexp.MustCompile(`table.([a-z]*) \[(.*)\]`)

	for _, table := range b.project.Tables {
		table.Permissions = []*TableFunctionPermision{}

		for _, function := range b.project.Functions {
			file, err := os.Open(filepath.Join(function.SourcePath, "main.go"))
			if err != nil {
				log.Fatal(err)
			}

			defer func() {
				if err := file.Close(); err != nil {
					log.Fatalf("Failed to close file: %v", err)
				}
			}()

			scanner := bufio.NewScanner(file)
			matches := [][]string{}
			for scanner.Scan() {
				line := scanner.Text()
				match := reg.FindStringSubmatch(line)
				if len(match) > 0 {
					matches = append(matches, match)
				}
			}

			if len(matches) == 0 {
				continue
			}

			for _, match := range matches {
				tname := match[1]
				perms := strings.Split(match[2], " ")
				if tname == table.Name {
					permission := &TableFunctionPermision{
						Function:   function,
						Permisions: perms,
					}

					table.Permissions = append(table.Permissions, permission)
				}
			}
		}
	}
}

func (b *projectBuilder) loadResources(apiFolder string) []*Resource {
	reg := regexp.MustCompile(`api.(GET|POST|PUT|OPTIONS|DELETE) (.*) \[([ 0-9]+)\]`)

	resources := []*Resource{}

	// Loop find resources
	for _, function := range b.apiFunctions[apiFolder] {

		file, err := os.Open(filepath.Join(function.SourcePath, "main.go"))
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				log.Fatalf("Failed to close file: %v", err)
			}
		}()

		scanner := bufio.NewScanner(file)
		matches := [][]string{}
		for scanner.Scan() {
			line := scanner.Text()
			match := reg.FindStringSubmatch(line)
			if len(match) > 0 {
				matches = append(matches, match)
			}
		}

		if len(matches) == 0 {
			continue
		}

		resource := &Resource{
			Name: matches[0][2],
			Path: matches[0][2],
		}

		found := false
		for _, r := range b.project.Resources {
			if r.Name == resource.Name {
				found = true
			}
		}
		if !found {
			b.project.Resources = append(b.project.Resources, resource)
			resources = append(resources, resource)
		}
	}

	// Loop find methods
	mapParents := map[string]string{}
	for _, function := range b.apiFunctions[apiFolder] {
		file, err := os.Open(filepath.Join(function.SourcePath, "main.go"))
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				log.Fatalf("Failed to close file: %v", err)
			}
		}()

		scanner := bufio.NewScanner(file)
		matches := [][]string{}
		for scanner.Scan() {
			line := scanner.Text()
			match := reg.FindStringSubmatch(line)
			if len(match) > 0 {
				matches = append(matches, match)
			}
		}

		if len(matches) == 0 {
			continue
		}

		for _, match := range matches {
			methodName := match[1]
			resourceName := match[2]

			if strings.Contains(methodName, "/") {
				parts := strings.Split(methodName, "/")
				mapParents[parts[0]] = parts[1]
				methodName = parts[1]
			}

			method := &Method{
				Function:  function,
				Method:    methodName,
				Responses: []string{},
			}

			resps := strings.Split(match[3], " ")
			method.Responses = append(method.Responses, resps...)

			for ir, r := range b.project.Resources {
				if r.Name == resourceName {
					b.project.Resources[ir].Methods = append(b.project.Resources[ir].Methods, method)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	for _, rp := range b.project.Resources {
		for ic, ac := range b.project.Resources {
			if strings.Contains(ac.Name, "/") {
				parts := strings.Split(ac.Name, "/")
				if len(parts) > 1 {
					pname := parts[0]
					actual := parts[1]
					if pname == rp.Name {
						bname := strings.Replace(ac.Name, "/", "_", -1)
						bname = strings.Replace(bname, "{", "", -1)
						bname = strings.Replace(bname, "}", "", -1)
						b.project.Resources[ic].Name = bname
						b.project.Resources[ic].Path = actual
						b.project.Resources[ic].Parent = rp
					}
				}
			}
		}
	}

	return resources
}

func (b *projectBuilder) loadFunctions() []*Function {
	functions := []*Function{}

	apiFolders := listFolders(b.config.FunctionsSourceDir)

	for _, apiFolder := range apiFolders {
		b.apiFunctions[apiFolder] = []*Function{}

		serviceDir := filepath.Join(b.config.FunctionsSourceDir, apiFolder)

		functionFolders := listFolders(serviceDir)
		for _, functionFolder := range functionFolders {
			function := &Function{
				Name:       functionFolder,
				SourcePath: filepath.Join(serviceDir, functionFolder),
			}
			functions = append(functions, function)
			b.apiFunctions[apiFolder] = append(b.apiFunctions[apiFolder], function)
		}
	}

	return functions
}

func (b *projectBuilder) loadTables() []*Table {
	tables := []*Table{}
	for _, tableName := range b.listTables() {
		table := &Table{
			Name: tableName,
		}
		tables = append(tables, table)

	}
	return tables
}

func (b *projectBuilder) listTables() []string {
	tables := []string{}

	apiFolders := listFolders(b.config.FunctionsSourceDir)

	for _, apiFolder := range apiFolders {
		jsonPath := filepath.Join(b.config.FunctionsSourceDir, apiFolder, "tables.json")

		jsonFile, err := os.Open(jsonPath)
		if err != nil {
			panic(err)
		}

		defer func() {
			if err := jsonFile.Close(); err != nil {
				log.Fatalf("Failed to close file: %v", err)
			}
		}()

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			panic(err)
		}

		var jtables []struct {
			Name string `json:"name"`
		}

		json.Unmarshal(byteValue, &jtables)
		for _, t := range jtables {
			tables = append(tables, t.Name)
		}
	}

	return tables
}
