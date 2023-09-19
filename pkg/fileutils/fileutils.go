package fileutils

import "os"

// Checks if fileA is newer than fileB by comparing ModTime.
func IsNewer(fileA, fileB string) (bool, error) {
	statA, err := os.Stat(os.ExpandEnv(fileA))
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	statB, err := os.Stat(os.ExpandEnv(fileB))
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	return statA.ModTime().After(statB.ModTime()), nil
}
