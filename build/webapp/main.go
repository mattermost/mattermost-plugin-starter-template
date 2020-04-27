// main exits with a non-0 status code if its necessary to run npm install
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) == 1 {
		help()
		return
	}

	if len(os.Args) != 3 {
		fmt.Println("Unexpected arguments.")
		help()
		os.Exit(1)
	}

	cmd := os.Args[1]
	webappPath := os.Args[2]
	switch cmd {
	case "check":
		matches, err := check(webappPath)
		if err != nil {
			fmt.Printf("Failed to check if .npminstall contains package.json hash: %s\n", err.Error())
			os.Exit(1)
		}
		if !matches {
			os.Exit(2)
		}

	case "update":
		err := update(webappPath)
		if err != nil {
			fmt.Printf("Failed to update .npminstall with hash of package.json: %s\n", err.Error())
			os.Exit(1)
		}

	default:
		fmt.Printf("Unexpected command %s\n\n", cmd)
		help()

		os.Exit(1)
	}
}

// help outputs instructions on how to use the command
func help() {
	fmt.Println("Usage:")
	fmt.Println("    webapp check <webapp path>")
	fmt.Println("    webapp update <webapp path>")
	fmt.Println()
}

// check checks if .npminstall contains the sha256 of package.json.
func check(webappPath string) (bool, error) {
	packageJSONPath := filepath.Join(webappPath, "package.json")
	packageJSONHash, err := hashFile(packageJSONPath)
	if err != nil {
		return false, errors.Wrap(err, "failed to hash")
	}

	npmInstallPath := filepath.Join(webappPath, ".npminstall")
	npmInstallHash, err := readHashFromFile(npmInstallPath)
	if err != nil {
		return false, errors.Wrap(err, "failed to read hash")
	}

	if len(npmInstallHash) == 0 {
		fmt.Printf("no previously recorded hash of %s (%x) in %s\n", packageJSONPath, packageJSONHash, npmInstallPath)
		return false, nil
	}

	if bytes.Equal(packageJSONHash, npmInstallHash) {
		return true, nil
	}

	fmt.Printf("hash of %s (%x) different from value recorded in %s (%x)\n", packageJSONPath, packageJSONHash, npmInstallPath, npmInstallHash)

	return false, nil
}

// update updates .npminstall with the sha256 of package.json.
func update(webappPath string) error {
	packageJSONPath := filepath.Join(webappPath, "package.json")
	packageJSONHash, err := hashFile(packageJSONPath)
	if err != nil {
		return errors.Wrap(err, "failed to hash")
	}

	npmInstallPath := filepath.Join(webappPath, ".npminstall")
	err = ioutil.WriteFile(npmInstallPath, []byte(hex.EncodeToString(packageJSONHash)), 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write hash to %s", npmInstallPath)
	}

	return nil
}

// hashFile computes the sha256 hash of the given file
func hashFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %s", filePath)
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return nil, errors.Wrapf(err, "failed to hash %s", filePath)
	}

	return hash.Sum(nil), nil
}

// readHashFromFile recovers a hash previously written to the given file
func readHashFromFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to open %s", filePath)
	}

	hashString, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %s", filePath)
	}

	if len(hashString) == 0 {
		return nil, nil
	}

	hash, err := hex.DecodeString(string(hashString))
	if err != nil {
		fmt.Printf("ignoring unexpected hash string in %s: %s\n", filePath, string(hashString))
		return nil, nil
	}

	return hash, nil
}
