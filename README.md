# ReExporter

ReExporter is a tool for generating code that exports symbols from sub-packages

## Usage

1. Install ReExpert by running
   `go install github.com/marvinpeter95/reexporter@main`
2. Place exported.yaml into a Go package which should export symbols from
   sub-package. See [example/](/example/).
3. Run `reexporter` from the root of your project.
