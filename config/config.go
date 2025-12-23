package config

import (
	"go/ast"
	"maps"
	"os"

	"github.com/goccy/go-yaml"
)

type ExportType string

const (
	ExportTypeType     ExportType = "type"
	ExportTypeVariable ExportType = "variable"
	ExportTypeConstant ExportType = "constant"
	ExportTypeFunction ExportType = "function"
)

// Config represents the overall configuration for the re-exporter.
type Config struct {
	Common  Export   `yaml:"common"`  // Common export configuration
	Exports []Export `yaml:"exports"` // List of export configurations
}

// Export represents the export configuration for a specific module.
type Export struct {
	Import  string            `yaml:"import"`  // Module import path
	Output  string            `yaml:"output"`  // Output file name
	Exclude Exclusion         `yaml:"exclude"` // Exclusion rules for re-exports
	Rename  map[string]string `yaml:"rename"`  // Rename symbol name during re-export
}

// Exclusion defines what kinds of symbols to exclude from re-exporting.
type Exclusion struct {
	Types     bool     `yaml:"types"`     // Do not export types
	Variables bool     `yaml:"variables"` // Do not export variables
	Constants bool     `yaml:"constants"` // Do not export constants
	Functions bool     `yaml:"functions"` // Do not export functions
	Names     []Filter `yaml:"names"`     // Do not export names matching these filters
	Files     []Filter `yaml:"files"`     // Do not export names from file matching these filters (file name only without extension)
}

// IncludeFile checks if the given file name is allowed based on the
// exclusion rules. It returns true if the file is included, false if excluded.
func (es *Export) IncludeFile(fileName string) bool {
	for _, f := range es.Exclude.Files {
		if f.Match(fileName) {
			return false
		}
	}
	return true
}

// ExportAs determines the export name for a given identifier based on the
// export configuration. It returns the new name and a boolean indicating
// whether the identifier should be exported.
func (es *Export) ExportAs(name *ast.Ident, exportType ExportType) (string, bool) {
	// Validate name and exportability
	if name == nil || name.Name == "" || !name.IsExported() {
		return "", false
	}

	// Check type-based exclusions
	if exportType == ExportTypeType && es.Exclude.Types ||
		exportType == ExportTypeVariable && es.Exclude.Variables ||
		exportType == ExportTypeConstant && es.Exclude.Constants ||
		exportType == ExportTypeFunction && es.Exclude.Functions {
		return name.Name, false
	}

	// Check exclude filters first
	for _, f := range es.Exclude.Names {
		if f.Match(name.Name) {
			return name.Name, false
		}
	}

	// Apply renaming if applicable
	newName := name.Name
	if renamed, ok := es.Rename[name.Name]; ok {
		newName = renamed
	}

	// Not matched by specified include filter
	return newName, true
}

// FromFile loads the exporter configuration from the given YAML file.
func FromFile(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	// Set defaults and merge common settings
	if config.Common.Output == "" {
		config.Common.Output = "exported.go"
	}

	for i := range config.Exports {
		es := &config.Exports[i]

		// Merge Exclude
		es.Exclude.Types = es.Exclude.Types || config.Common.Exclude.Types
		es.Exclude.Variables = es.Exclude.Variables || config.Common.Exclude.Variables
		es.Exclude.Constants = es.Exclude.Constants || config.Common.Exclude.Constants
		es.Exclude.Functions = es.Exclude.Functions || config.Common.Exclude.Functions
		es.Exclude.Names = append(es.Exclude.Names, config.Common.Exclude.Names...)
		es.Exclude.Files = append(es.Exclude.Files, config.Common.Exclude.Files...)

		if es.Output == "" {
			es.Output = config.Common.Output
		}

		// Merge Rename
		maps.Copy(es.Rename, config.Common.Rename)
	}

	return &config, nil
}
