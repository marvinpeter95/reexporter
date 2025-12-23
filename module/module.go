package module

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// ErrGoModNotFound is returned when no go.mod file is found in the directory hierarchy.
var ErrGoModNotFound = errors.New("go.mod not found")

var modCache map[string]*modfile.File

func GetModuleFor(pkgDir string) (string, *modfile.File, error) {
	// Initialize the modCache map if it hasn't been already
	if modCache == nil {
		modCache = make(map[string]*modfile.File)
	}

	// Get the absolute path of the package directory
	modDir, err := filepath.Abs(pkgDir)
	if err != nil {
		return "", nil, err
	}

	// Check if the module file is already cached
	if mod, ok := modCache[modDir]; ok {
		return modDir, mod, nil
	}

	// Traverse up the directory tree to find the nearest go.mod file
	for {
		goModFile := filepath.Join(modDir, "go.mod")
		_, err := os.Stat(goModFile)
		// If the go.mod file does not exist, move to the parent directory
		if os.IsNotExist(err) {
			parent := filepath.Dir(modDir)
			if parent == modDir {
				break
			}
			modDir = parent
			continue
		} else if err != nil {
			return "", nil, err
		}

		modCache[modDir], err = loadGoMod(goModFile)
		if err != nil {
			return "", nil, err
		}

		return modDir, modCache[modDir], nil
	}

	return "", nil, ErrGoModNotFound
}

// LoadGoMod loads and parses a go.mod file from the specified file path.
func loadGoMod(filePath string) (*modfile.File, error) {
	goModData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	mod, err := modfile.Parse(filePath, goModData, nil)
	if err != nil {
		return nil, err
	}

	return mod, nil
}
