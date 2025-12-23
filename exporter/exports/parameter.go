package exports

import (
	"go/ast"
	"go/types"
	"strings"
)

// Parameter represents a function parameter or return value.
type Parameter struct {
	Name     string // The name of the parameter (optional).
	Type     string // The type of the parameter.
	Variadic bool   // Indicates if the parameter is variadic.
}

// Variable returns the parameter as it would appear in a function call.
func (p Parameter) Variable() string {
	if p.Variadic {
		return p.Name + "..."
	}
	return p.Name
}

// Parameter returns the parameter as it would appear in a function signature.
func (p Parameter) Parameter() string {
	sb := &strings.Builder{}
	if p.Name != "" {
		sb.WriteString(p.Name)
		sb.WriteString(" ")
	}
	if p.Variadic {
		sb.WriteString("...")
	}
	sb.WriteString(p.Type)
	return sb.String()
}

// String returns the string representation of the parameter.
func (p Parameter) String() string {
	return p.Parameter()
}

// FunctionSignature represents the signature of a function, including type parameters, parameters, and results.
type FunctionSignature struct {
	Types      []Parameter
	Parameters []Parameter
	Results    []Parameter
}

// ParseFunctionSignature parses the function signature from an AST function declaration.
func ParseFunctionSignature(decl *ast.FuncDecl) FunctionSignature {
	sig := FunctionSignature{
		Types:      []Parameter{},
		Parameters: []Parameter{},
		Results:    []Parameter{},
	}

	// Handle type parameters
	if decl.Type.TypeParams != nil {
		for _, field := range decl.Type.TypeParams.List {
			sig.Types = append(sig.Types, parameterFromField(field)...)
		}
	}

	// Handle parameters
	if decl.Type.Params != nil {
		for _, field := range decl.Type.Params.List {
			sig.Parameters = append(sig.Parameters, parameterFromField(field)...)
		}
	}

	// Handle results
	if decl.Type.Results != nil {
		for _, field := range decl.Type.Results.List {
			sig.Results = append(sig.Results, parameterFromField(field)...)
		}
	}

	return sig
}

// parameterFromField creates a Parameter from an AST field.
func parameterFromField(field *ast.Field) []Parameter {
	names := field.Names

	// Handle unnamed parameters
	if len(field.Names) == 0 {
		names = append(names, &ast.Ident{Name: ""})
	}

	ps := make([]Parameter, len(names))
	for i, nameIdent := range names {
		p := Parameter{Name: nameIdent.Name}
		p.Type, p.Variadic = strings.CutPrefix(types.ExprString(field.Type), "...")

		ps[i] = p
	}

	return ps
}
