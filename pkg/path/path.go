package path

import (
	"path/filepath"
	"runtime"
)

func RootDir(paths ...string) string {
	_, filename, _, _ := runtime.Caller(0) // Path for this file
    currentDir := filepath.Dir(filename)      // Directory of the file
    root := filepath.Join(currentDir, "..", "..") // Go up to the root of the project

    if len(paths) > 0 {
        root = filepath.Join(root, filepath.Join(paths...))
    }

    return root
}
