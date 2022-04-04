//go:build !prod
// +build !prod

package rebuild

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	force              bool
	ErrSourceModified  = errors.New("source modified")
	ErrCheckerFailed   = errors.New("failed to verify source timestamps")
	ErrSourcesNotFound = errors.New("source files not found")
)

func Run(force bool) {
	path, err := RebuildChecker()
	if err == nil {
		return
	}

	if errors.Is(err, ErrSourceModified) {
		fmt.Fprintf(os.Stderr, "*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*\n")
		fmt.Fprintf(os.Stderr, "*~*                             *~*\n")
		if force {
			fmt.Fprintf(os.Stderr, "*~*   ERROR: sources modified   *~*\n")
		} else {
			fmt.Fprintf(os.Stderr, "*~*  WARNING: sources modified  *~*\n")
		}
		fmt.Fprintf(os.Stderr, "*~*      Please rebuild me!     *~*\n")
		fmt.Fprintf(os.Stderr, "*~*                             *~*\n")
		fmt.Fprintf(os.Stderr, "*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*\n")
		fmt.Fprintf(os.Stderr, "Modified source: %s\n\n", path)

		if force {
			os.Exit(127)
		}
	} else {
		fmt.Fprintf(os.Stderr, "WARNING: rebuild checker failed: %v\n", err)
	}
}

func RebuildChecker() (string, error) {
	myName, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("%w: getting path to current executable: %v", ErrCheckerFailed, err)
	}

	mySt, err := os.Stat(myName)
	if err != nil {
		return "", fmt.Errorf("%w: stat() failed: %v", ErrCheckerFailed, err)
	}

	myModTime := mySt.ModTime()

	myDir := filepath.Dir(myName)
	glob := filepath.Join(myDir, "*.go")

	gofiles, err := filepath.Glob(glob)
	if err != nil {
		return "", fmt.Errorf("%w: listing directory: %v", ErrCheckerFailed, err)
	}

	if len(gofiles) == 0 {
		return "", fmt.Errorf("%w: source files for %s not found in %s", ErrCheckerFailed, myName, myDir)
	}

	rootDir := myDir
	prevRootDir := ""

	// find the root of the source tree
	for rootDir != prevRootDir {
		prevRootDir = rootDir

		glob := filepath.Join(rootDir, "go.mod")
		goMod, err := filepath.Glob(glob)
		if err != nil {
			return "", fmt.Errorf("%w: listing directory: %v", ErrCheckerFailed, err)
		}

		if len(goMod) == 1 {
			var path string

			err = filepath.WalkDir(rootDir, fileCheckFunc(myModTime, &path))
			if err == nil {
				return "", nil
			}

			if errors.Is(err, ErrSourceModified) {
				return path, err
			}

			return "", fmt.Errorf("%w: %v", ErrCheckerFailed, err)
		}

		rootDir = filepath.Dir(rootDir)
	}

	return "", fmt.Errorf("%w: sources not found in %s", ErrCheckerFailed, myDir)
}

func fileCheckFunc(t time.Time, p *string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.ModTime().After(t) {
			*p = path

			return ErrSourceModified
		}

		return nil
	}
}
