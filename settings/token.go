// Copyright (c) 2015-2018 TIBCO Software Inc.
// All Rights Reserved

package settings

import "net/http"

const (
	TOKEN_FILE_NAME              = "token"
	ACCESS_TOKEN_KEY_NAME string = "AccessToken"
)

// The following Token struct present 'token' file which locates in '~/.tibcli' folder and it contains TA accessToken
type Token struct {
	// serializable fields
	AccessToken *http.Cookie // AccessToken keeps TA accessToken after TA login passed

	// non-serializable (i.e. private) fields
	*settingsFile // base type, containing all logic for serialization & deserialization
}

// NewToken creates a new token object.
func NewToken() (*Token, error) {
	settingsFile, err := newSettingsFile(TOKEN_FILE_NAME)
	if err != nil {
		return nil, err
	}
	t := &Token{settingsFile: settingsFile}
	return t, nil
}

// Read loads Task from disk.
// If the corresponding disk file is empty, all public fields will be empty.
func (t *Token) Read(unobfuscateValue bool) (err error) {
	if err := t.read(t); err != nil {
		return err
	} else {
		if unobfuscateValue && t.AccessToken != nil {
			t.AccessToken.Value = unobfuscate(t.AccessToken.Value)
		}
	}
	return nil
}

// Write saves the current object to disk.
func (t *Token) Write(obfuscateValue bool) error {

	if obfuscateValue && t.AccessToken != nil {
		// create a copy of s, to preserve the rest of the fields
		var s_obfus *Token = new(Token)
		*s_obfus = *t
		newCookie := &http.Cookie{
			Name:    s_obfus.AccessToken.Name,
			Value:   obfuscate(s_obfus.AccessToken.Value),
			Path:    s_obfus.AccessToken.Path,
			Domain:  s_obfus.AccessToken.Domain,
			Expires: s_obfus.AccessToken.Expires,
		}
		s_obfus.AccessToken = newCookie //Override the cookie
		return t.write(s_obfus)
	} else {
		return t.write(t)
	}
}

// Deletes this token file from disk.
func (t *Token) Delete() error {
	return t.deleteFile()
}
