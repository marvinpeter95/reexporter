package exporter

import (
	_ "embed"
	"errors"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/marvinpeter95/reexporter/config"
	"github.com/marvinpeter95/reexporter/exporter/exports"
	"golang.org/x/tools/go/packages"
)

var ErrLoadingPackages = errors.New("failed to load packages")

// Exporter represents the code exporter.
type Exporter struct {
	Exports []config.Export  // The export configurations.
	Dir     string           // The dictory of the main module where the go.mod is located.
	PkgName string           // The package name for the generated code.
	data    *exports.Exports // Holds the collected export data.
	fset    *token.FileSet   // Keep track of positions for file-based exclusion.
}

// New creates a new Exporter with the given configuration.
func New(exports []config.Export, dir string, pkgName string) *Exporter {
	return &Exporter{Exports: exports, Dir: dir, PkgName: pkgName}
}

// Generate generates the exported code based on the configuration.
func (e *Exporter) Generate() (string, error) {
	e.data = exports.New(filepath.Base(e.PkgName))
	e.fset = token.NewFileSet()

	for _, export := range e.Exports {
		if err := e.processExport(export); err != nil {
			return "", err
		}
	}

	// Render the template with the collected data.
	codeStr, err := renderTemplate(e.data)
	if err != nil {
		return "", err
	}

	// Format the generated code.
	formatted, err := formatCode(codeStr)
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

// processExport processes a single export configuration and updates the ExportData accordingly.
func (e *Exporter) processExport(export config.Export) error {
	// Resolve relative imports.
	if sub, ok := strings.CutPrefix(export.Import, "./"); ok {
		export.Import = filepath.Join(e.PkgName, sub)
	}

	// Always add the main import.
	e.data.AddImport(export.Import)

	cfg := packages.Config{
		Fset: e.fset,
		Dir:  e.Dir,
		Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedImports,
	}

	// Load the imported package
	pkgs, err := packages.Load(&cfg, export.Import)
	if err != nil {
		return err
	}

	// Check for errors while loading packages
	if packages.PrintErrors(pkgs) > 0 {
		return ErrLoadingPackages
	}

	// Process each package
	for _, pkg := range pkgs {
		// Add necessary imports from the package
		for _, imp := range pkg.Imports {
			// Only add imports that are not part of the standard library, since those may
			// may not be resolvable via the formatter.
			if strings.Contains(imp.ID, ".") {
				e.data.AddImport(imp.ID)
			}
		}

		// Inspect the AST of each file in the package
		for _, fileAst := range pkg.Syntax {
			ast.Inspect(fileAst, func(n ast.Node) bool {
				return e.inspectAST(pkg, &export, n)
			})
		}
	}

	return nil
}

// inspectAST inspects the AST nodes and collects exportable entities based on the export configuration.
func (e *Exporter) inspectAST(pkg *packages.Package, export *config.Export, n ast.Node) bool {
	if n == nil {
		return true
	}
	fn := e.fset.File(n.Pos()).Name()
	fn = strings.TrimSuffix(filepath.Base(fn), ".go")

	if !export.IncludeFile(fn) {
		return true
	}

	switch n := n.(type) {
	case *ast.GenDecl:
		for _, spec := range n.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				if name, ok := export.ExportAs(s.Name, config.ExportTypeType); ok {
					e.data.AddType(name, s.Name.Name, filepath.Base(pkg.ID), exports.ParseComment(n.Doc, s.Comment))
				}
			case *ast.ValueSpec:
				for _, nameIdent := range s.Names {
					exportType := config.ExportTypeVariable
					if n.Tok == token.CONST {
						exportType = config.ExportTypeConstant
					}
					if name, ok := export.ExportAs(nameIdent, exportType); ok {
						if exportType == config.ExportTypeVariable {
							e.data.AddVariable(name, nameIdent.Name, filepath.Base(pkg.ID), exports.ParseComment(n.Doc, s.Comment))
						} else if exportType == config.ExportTypeConstant {
							e.data.AddConstant(name, nameIdent.Name, filepath.Base(pkg.ID), exports.ParseComment(n.Doc, s.Comment))
						}
					}
				}
			}
		}
	case *ast.FuncDecl:
		if name, ok := export.ExportAs(n.Name, config.ExportTypeFunction); ok && n.Recv == nil {
			e.data.AddFunction(name, n.Name.Name, filepath.Base(pkg.ID), exports.ParseComment(n.Doc, nil), exports.ParseFunctionSignature(n))
		}
	default:
		return true
	}
	return false
}
