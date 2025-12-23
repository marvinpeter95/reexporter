package exports

import (
	"slices"
	"strings"
)

// Exports holds all the collected export data.
type Exports struct {
	Pkg        string              // Package name
	Imports    []string            // List of import paths
	Types      []Export            // List of type exports
	Variables  []Export            // List of variable exports
	Constants  []Export            // List of constant exports
	Functions  []FunctionExport    // List of function exports
	importsSet map[string]struct{} // Set to track unique imports
}

// New creates a new ExportData instance for the given package.
func New(pkg string) *Exports {
	return &Exports{
		Pkg:        pkg,
		Imports:    []string{},
		Types:      []Export{},
		Variables:  []Export{},
		Constants:  []Export{},
		Functions:  []FunctionExport{},
		importsSet: make(map[string]struct{}),
	}
}

// AddImport adds a new import path if it doesn't already exist.
func (td *Exports) AddImport(importPath string) {
	if _, exists := td.importsSet[importPath]; !exists {
		td.importsSet[importPath] = struct{}{}
		td.Imports = append(td.Imports, importPath)
	}
}

// AddType adds a new type export.
func (td *Exports) AddType(exportName string, name string, pkg string, c Comment) {
	td.Types = insertSortedExport(td.Types, Export{
		ExportName: exportName,
		Name:       name,
		Package:    pkg,
		Comment:    c,
	})
}

// AddVariable adds a new variable export.
func (td *Exports) AddVariable(exportName string, name string, pkg string, c Comment) {
	td.Variables = insertSortedExport(td.Variables, Export{
		ExportName: exportName,
		Name:       name,
		Package:    pkg,
		Comment:    c,
	})
}

// AddConstant adds a new constant export.
func (td *Exports) AddConstant(exportName string, name string, pkg string, c Comment) {
	td.Constants = insertSortedExport(td.Constants, Export{
		ExportName: exportName,
		Name:       name,
		Package:    pkg,
		Comment:    c,
	})
}

// AddFunction adds a new function export.
func (td *Exports) AddFunction(exportName string, name string, pkg string, c Comment, sig FunctionSignature) {
	td.Functions = append(td.Functions, FunctionExport{
		Export: Export{
			ExportName: exportName,
			Name:       name,
			Package:    pkg,
			Comment:    c,
		},
		Signature: sig,
	})
}

// insertSortedExport inserts an Export into a sorted slice while maintaining order.
func insertSortedExport(ts []Export, t Export) []Export {
	i, _ := slices.BinarySearchFunc(ts, t, func(a, b Export) int {
		return strings.Compare(a.ExportName, b.ExportName)
	})
	return slices.Insert(ts, i, t)
}
