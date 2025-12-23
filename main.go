package main

import (
	_ "embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/marvinpeter95/reexporter/config"
	"github.com/marvinpeter95/reexporter/exporter"
	"github.com/marvinpeter95/reexporter/module"
)

func main() {
	// Start looking for exported.yaml files from the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Walk the directory tree to find exported.yaml files
	err = filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
		// If there's an error or it's not the exported.yaml file, skip it
		if err != nil || d.IsDir() || d.Name() != "exported.yaml" {
			return err
		}

		// Load the configuration from the exported.yaml file
		config, err := config.FromFile(path)
		if err != nil {
			return err
		}

		// Get the module information closest to the current directory
		baseModDir, mod, err := module.GetModuleFor(filepath.Dir(path))
		if err != nil {
			return err
		}

		// Determine the package path relative to the module
		rel, err := filepath.Rel(baseModDir, filepath.Dir(path))
		if err != nil {
			return err
		}

		// Join the module path with the relative path to get the full package path
		pkgPath := filepath.Join(mod.Module.Mod.Path, rel)
		println(pkgPath)

		// Create a new exporter and generate the code
		exporter := exporter.New(config.Exports, baseModDir, pkgPath)
		code, err := exporter.Generate()
		if err != nil {
			return err
		}

		// Write the generated code next to the exported.yaml file
		outputName := config.Common.Output
		if baseName, ok := strings.CutPrefix(outputName, "__"); ok {
			outputName = filepath.Base(pkgPath) + baseName
		}
		return os.WriteFile(
			filepath.Join(filepath.Dir(path), outputName),
			[]byte(code),
			0o644,
		)
	})
	if err != nil {
		panic(err)
	}
}
