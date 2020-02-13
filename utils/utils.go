package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Morphyni/tas-cli/consts"
	"github.com/Morphyni/tas-cli/settings"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
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

// IncorrectUsageError is a custom error representing an incorrect command usage
type IncorrectUsageError struct {
	Context *cli.Context
	Msg     string
}

func (e *IncorrectUsageError) Error() string {
	return e.Msg
}

// CheckError is a generic error checker. If the supplied error nil, this is a no-op. If
// the error is an IncorrectUsageError, it displays the help for the command, else it displays
// the error and exit the current application
func CheckError(err error) {
	if err != nil {
		switch e := err.(type) {
		case *IncorrectUsageError:
			fmt.Printf("Incorrect Usage: %s\n\n", e)
			cli.ShowSubcommandHelp(e.Context)
		default:
			fmt.Printf("Error: %v\n", e)
		}
		os.Exit(1)
	}
}

// GetEnvParam get env param value
func GetEnvParam(name string) string {
	value := os.Getenv(name)
	if value != "" {
		log.Debugf("Env param '%s' is found, it has value: %s", name, value)
	} else {
		log.Debugf("Env param '%s' is not found", name)
	}
	return value
}

// LoadToken loads user TA AccessToken
func LoadToken(unobfuscate bool) (*settings.Token, error) {
	token, err := settings.NewToken()
	if err != nil {
		return nil, err
	}
	err = token.Read(unobfuscate)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// LoadSettings loads user's all token/profile/session from ~/.tibcli
func LoadSettings() (*settings.Profile, *settings.Session, *settings.Token, error) {
	profile, err := LoadProfile()
	if err != nil {
		return nil, nil, nil, err
	}
	session, err := LoadSession(consts.OBFUSCATE_COOKIE_VALUE)
	if err != nil {
		return nil, nil, nil, err
	}
	token, err := LoadToken(consts.OBFUSCATE_COOKIE_VALUE)
	if err != nil {
		return nil, nil, nil, err
	}
	return profile, session, token, err
}

// GetOrgAndRegion gets org and region info from session
func GetOrgAndRegion() (string, string, error) {
	var org string
	var region string
	session, err := LoadSession(consts.OBFUSCATE_COOKIE_VALUE)
	if len(session.OrgName) > 0 {
		org = session.OrgName
	}

	if len(session.OrgDisplayName) > 0 {
		orgDisplayName := session.OrgDisplayName
		if strings.HasPrefix(orgDisplayName, org) {
			s := strings.Split(orgDisplayName, org)
			region = strings.TrimSpace(s[len(s)-1])
		}
	}
	return org, region, err
}

// LoadSession loads user Session
func LoadSession(unobfuscate bool) (*settings.Session, error) {
	session, err := settings.NewSession()
	if err != nil {
		return nil, err
	}
	err = session.Read(unobfuscate)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// GetDomainURL get updated domain URL.
// session.DomainUrl contains updated Domain Server Host URL, and use placeholder value if session.DomainUrl is empty.
func GetDomainURL() (string, error) {
	session, err := LoadSession(consts.OBFUSCATE_COOKIE_VALUE)
	if err != nil {
		return "", err
	}
	if session.DomainUrl != "" {
		return session.DomainUrl, nil
	}
	err, domainUrl := settings.GetPlaceHolderValue(settings.DOMAIN_SERVER_HOST_PLACEHOLDER)
	if err != nil {
		return "", err
	}
	return domainUrl, nil
}

// GetIDMConnectURL() get idm connect URL
func GetIDMConnectURL() (string, error) {
	err, idmConnectUrl := settings.GetPlaceHolderValue(settings.IDENTITY_MANAGEMENT_SERVER_HOST_PLACEHOLDER)
	if err != nil {
		return "", err
	}
	return idmConnectUrl, nil
}

// PromptForUser interactively prompts for a username input
func PromptForUser(user string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username(" + user + "): ")
	inputUsr, _ := reader.ReadString('\n')
	inputUsr = strings.TrimSpace(inputUsr)
	if len(inputUsr) == 0 {
		inputUsr = user
	}
	return inputUsr
}

// PromptForPassword interactively prompts for a passwrod input
func PromptForPassword() string {
	fmt.Print("Password: ")
	pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err == io.EOF {
		log.Debugf("ReadPassword error: %s", err.Error())
		os.Exit(1)
	}
	CheckError(err)
	fmt.Println()
	return string(pwd)
}
