package infra

import (
	"path/filepath"
	"runtime"
)

const (
	functionsSourceDir = "./functions"
	functionsBinaryDir = "./data/functions"
)

type Config struct {
	BaseDir            string `json:"base_dir"`
	ServicesDir        string `json:"services_dir"`
	FunctionsSourceDir string `json:"functions_dir"`
	FunctionsBinaryDir string `json:"function_binary_dir"`
}

func NewConfig(basePath string) Config {
	return Config{
		BaseDir:            basePath,
		FunctionsSourceDir: filepath.Join(basePath, functionsSourceDir),
		FunctionsBinaryDir: filepath.Join(basePath, functionsBinaryDir),
	}
}

var TestConfig Config

func init() {
	_, filename, _, _ := runtime.Caller(0)
	projectData := filepath.Dir(filename)
	TestConfig = NewConfig(filepath.Join(projectData, "testdata"))
}
