package exports

// Export represents a single exportable entity.
type Export struct {
	ExportName string  // The name to be used in the export.
	Name       string  // The original name in the source package.
	Package    string  // The package from which the entity is exported.
	Comment    Comment // Associated documentation and comments.
}

// FunctionExport represents an exported function with its signature.
type FunctionExport struct {
	Export // Base export information.

	Signature FunctionSignature // The function signature details.
}
