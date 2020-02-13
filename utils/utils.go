package utils

import (
	"github.com/Morphyni/tas-cli/settings"
	log "github.com/sirupsen/logrus"
)

// GetUserEmail gets user email:
// 1. get user email from profile
// 2. if profile doesn't have user email, then get that from placeholder
func GetUserEmail() (string, error) {
	profile, err := LoadProfile()
	if err != nil {
		log.Debug(err.Error())
		return "", err
	}

	if len(profile.UserEmail) > 0 {
		return profile.UserEmail, nil
	}

	err, userEmail := settings.GetPlaceHolderValue(settings.USERNAME_PLACEHOLDER)
	if err != nil {
		log.Debug(err.Error())
		return "", err
	}

	return userEmail, nil
}

// LoadProfile loads user Profile
func LoadProfile() (*settings.Profile, error) {
	profile, err := settings.NewProfile()
	if err != nil {
		return nil, err
	}
	err = profile.Read()
	if err != nil {
		return nil, err
	}
	return profile, nil
}
