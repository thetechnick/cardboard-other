package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"cardboard.package-operator.run/pkg/fileutils"
)

type DependencyDirectory string

// Returns the /bin directory containing the dependency binaries.
func (d DependencyDirectory) Bin() string {
	return filepath.Join(string(d), "bin")
}

// Go install a dependency into the dependency directory
func (d DependencyDirectory) GoInstall(tool, packageURl, version string) error {
	if err := os.MkdirAll(string(d), os.ModePerm); err != nil {
		return fmt.Errorf("create dependency dir: %w", err)
	}

	needsRebuild, err := d.NeedsRebuild(tool, version)
	if err != nil {
		return err
	}
	if !needsRebuild {
		return nil
	}

	url := packageURl + "@" + version

	cmd := exec.Command("go", "install", url)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOBIN="+d.Bin())

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("install %s: %s: %w", url, string(out), err)
	}
	return nil
}

// Checks if a tool in the dependency directory needs to be rebuild.
func (d DependencyDirectory) NeedsRebuild(tool, version string) (needsRebuild bool, err error) {
	versionFile := filepath.Join(string(d), "versions", tool, version)
	if err := EnsureFile(versionFile); err != nil {
		return false, fmt.Errorf("ensure file: %w", err)
	}

	// Checks "tool" binary file modification date against version file.
	// If the version file is newer, tool is of the wrong version.
	rebuild, err := fileutils.IsNewer(versionFile, filepath.Join(d.Bin(), tool))
	if err != nil {
		return false, fmt.Errorf("rebuild check: %w", err)
	}

	return rebuild, nil
}

// Ensure a file and it's file path exist.
func EnsureFile(file string) error {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		f, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("creating file %s: %w", file, err)
		}
		defer f.Close()
		return nil
	}
	if err != nil {
		return fmt.Errorf("checking file %s: %w", file, err)
	}
	return nil
}
