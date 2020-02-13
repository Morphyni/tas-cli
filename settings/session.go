package settings

import (
	"errors"
	"net/http"

	"github.com/Morphyni/tas-cli/types"
	log "github.com/sirupsen/logrus"
)

// Session type represents the user session.
type Session struct {
	// serializable fields

	// we serialize cookies ONLY, and create cookieJar on the time of using it (like create the cookieJar for httpClient which is used in http request)
	Cookies []*http.Cookie `json:"cookies"`

	// following fields come from IDM login response
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
	UserId    string `json:"userId"`
	OrgName   string `json:"orgName"`
	TS        int    `json:"ts"`

	// following fields come from IDM login response when user belongs to multiple org/subscriptions.
	DomainUrl      string           `json:"domainUrl"`
	OrgDisplayName string           `json:"orgDisplayName"`
	OrgList        []types.OrgEntry `json:"orgList"`
	SubscriptionId string           `json:"subscriptionId"`

	// following fields come from DomainServer GetDefaultSandbox response
	DefaultSandboxName           string `json:"defaultSandboxName"`
	DefaultSandboxOrganizationId string `json:"defaultSandboxOrgId"`

	//Sandboxes      map[string]string // list of all sandboxes, retrieved from the login response

	// non-serializable (i.e. private) fields
	*settingsFile // base type, containing all logic for serialization & deserialization
}

// NewSession creates a new Session object.
func NewSession() (*Session, error) {
	settingsFile, err := newSettingsFile(SESSION_FILENAME)
	if err != nil {
		return nil, err
	}
	s := &Session{Cookies: []*http.Cookie{}, settingsFile: settingsFile}
	return s, nil
}

// Read loads Session from disk, unobfuscating cookie's value if present and so indicated.
// If the corresponding disk file is empty, all public fields will be empty.
func (s *Session) Read(unobfuscateValue bool) (err error) {
	if err := s.read(s); err != nil {
		return err
	} else {
		if unobfuscateValue && s.Cookies != nil {
			for _, cookie := range s.Cookies {
				cookie.Value = unobfuscate(cookie.Value)
			}
		}
	}
	return nil
}

// Write the current object to disk, optionally obfuscating the cookie's value.
func (s *Session) Write(obfuscateValue bool) error {

	if obfuscateValue && s.Cookies != nil {
		// create a copy of s, to preserve the rest of the fields
		var s_obfus *Session = new(Session)
		*s_obfus = *s
		var newCookies []*http.Cookie
		for _, cookie := range s.Cookies {
			// override the Cookie value in the copy with an obfuscated value
			newCookie := &http.Cookie{
				Name:    cookie.Name,
				Value:   obfuscate(cookie.Value),
				Path:    cookie.Path,
				Domain:  cookie.Domain,
				Expires: cookie.Expires,
			}
			newCookies = append(newCookies, newCookie)
		}
		s_obfus.Cookies = newCookies //Override the cookieJar
		return s.write(s_obfus)
	} else {
		return s.write(s)
	}
}

// Delete deletes the session file from disk
func (s *Session) Delete() error {
	return s.deleteFile()
}

// UpdateCookies will update session existing cookies with new values, optionally obfuscating the cookie's value.
func (s *Session) UpdateCookies(newCookies []*http.Cookie, obfuscateValue bool) error {

	log.Debugf("[UpdateCookies] obfuscate %+v, \n newCookies: %+v, \n s.Cookies(old Cookies): %+v", obfuscateValue, newCookies, s.Cookies)

	if newCookies == nil || len(newCookies) == 0 {
		log.Errorf("[UpdateCookies] new cookies value is empty: %+v", newCookies)
		return errors.New("[UpdateCookies] The given new cookes are empty")
	}
	if s.Cookies == nil || len(s.Cookies) == 0 { //first time set Cookies
		s.Cookies = newCookies
	} else {
		for _, newCookie := range newCookies {
			for i, oldCookie := range s.Cookies {
				if oldCookie.Name == newCookie.Name && oldCookie.Domain == newCookie.Domain && oldCookie.Path == newCookie.Path {
					log.Debugf("[UpdateCookies] Cookie '%s' get refreshed.", newCookie.Name)
					s.Cookies[i] = newCookie
				}
			}
		}
	}
	log.Debugf("Session cookies refreshed, new values: %+v", s.Cookies)

	s.Write(obfuscateValue)

	return nil
}