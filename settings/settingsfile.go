// Copyright (c) 2015-2017 TIBCO Software Inc.
// All Rights Reserved

// Package settings contains code to serialize & deserialize user settings to disk into the current user profile
package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
)

const (
	SETTINGS_DIR = ".tibcli"
	//Prefix put in front of obfuscated strings indicating the algorithm
	OBFUS_PREFIX string = "obfus_v1."
)

// Struct settingsFile is the base type for all files in the storage.  It contains methods for
// reading/writing the profile to disk. Any setting file implementation should extend this.
type settingsFile struct {
	filePath string
}

// newSettingsFile creates a new settingsFile object
func newSettingsFile(filename string) (*settingsFile, error) {
	sf := &settingsFile{}

	// compute profile filepath
	filePath, err := sf.getFilePath(filename)
	if err != nil {
		return nil, err
	}
	sf.filePath = filePath
	return sf, nil
}

// getSettingsDir returns the full path to the settings directory
func (sf *settingsFile) getSettingsDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(u.HomeDir, SETTINGS_DIR), nil
}

// getFilePath returns the full path to the given settings filename
// filename is the name of the file represented by this object (e.g. "profile", "settings", etc.)
func (sf *settingsFile) getFilePath(filename string) (string, error) {
	settingsDir, err := sf.getSettingsDir()
	if err != nil {
		return "", err
	}

	return path.Join(settingsDir, filename), nil
}

// fileExists returns true if the specifield file exists
func (sf *settingsFile) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// createSettingsDir create the settings directory if it does not exist. If
// it exists, it is a no-op
func (sf *settingsFile) createSettingsDir() error {
	settingsDir, err := sf.getSettingsDir()
	if err != nil {
		return err
	}

	if !sf.fileExists(settingsDir) {
		perm := os.FileMode(0700) // drwx------
		return os.MkdirAll(settingsDir, perm)
	}
	return nil
}

// write serializes the supplied object to the settings file on disk
// in is the struct that is to be serialized
func (sf *settingsFile) write(in interface{}) error {
	// marshal the current receiver to a slice of bytes
	bytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	// create directory (if needed)
	err = sf.createSettingsDir()
	if err != nil {
		return err
	}

	// create file
	perm := os.FileMode(0600) // -rw-------
	if err = ioutil.WriteFile(sf.filePath, bytes, perm); err != nil {
		return err
	}
	return nil
}

// read deserializes the content of the settings file from disk into the supplied object
// out is the struct where the file will be de-serialized into
func (sf *settingsFile) read(out interface{}) error {
	// read content from file (if exists)
	if sf.fileExists(sf.filePath) {
		bytes, err := ioutil.ReadFile(sf.filePath)
		if err != nil {
			return err
		}

		// unmarshal the file content to the current receiver
		if err = json.Unmarshal(bytes, out); err != nil {
			return err
		}
	}
	return nil
}

// deleteFile deletes settingsFile if it exists
func (sf *settingsFile) deleteFile() error {
	if sf.fileExists(sf.filePath) {
		return os.Remove(sf.filePath)
	}
	return nil
}

func obfuscate(s string) string {
	return OBFUS_PREFIX + rot13(s)
}

func rot13(s string) string {
	var o string = ""
	for _, ch := range s {
		if ch <= 'Z' && ch >= 'A' {
			o += string(rune(ch) + 32)
		} else if ch <= 'z' && ch >= 'a' {
			o += string(ch - 32)
		} else {
			o += string(ch)
		}
	}
	return o
}

func unobfuscate(obf string) string {
	if strings.HasPrefix(obf, OBFUS_PREFIX) {
		return rot13(obf[len(OBFUS_PREFIX):]) //we're doing ROT13; just strip the versioning prefix
	} else {
		return obf
	}
}