package settings

import "github.com/Morphyni/tas-cli/consts"

const (
	PROFILE_FILENAME string = "profile"
)

// Profile type represents the user profile.
type Profile struct {
	// serializable fields
	Version       string `json:"tibcliVersion"` // keeps tibcli version
	IDMConnectURL string `json:"idmConnectUrl"` // IDM server url
	UserEmail     string `json:"userEmail"`     // user email
	KnownRegion   string `json:"knownRegion"`   // known region

	// non-serializable (i.e. private) fields
	*settingsFile // base type, containing all logic for serialization & deserialization
}

// NewProfile creates a new Profile object.
func NewProfile() (*Profile, error) {
	settingsFile, err := newSettingsFile(PROFILE_FILENAME)
	if err != nil {
		return nil, err
	}
	p := &Profile{Version: consts.CLI_VERSION, settingsFile: settingsFile}
	return p, nil
}

// Read loads Profile from disk.
// If the corresponding disk file is empty, all public fields will be empty.
func (p *Profile) Read() error {
	return p.read(p)
}

// Write saves the current profile object to disk.
func (p *Profile) Write() error {
	return p.write(p)
}

// Deletes this profile from disk.
func (p *Profile) Delete() error {
	return p.deleteFile()
}
