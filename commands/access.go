package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Morphyni/tas-cli/consts"
	"github.com/Morphyni/tas-cli/settings"
	"github.com/Morphyni/tas-cli/types"
	"github.com/Morphyni/tas-cli/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func Login(c *cli.Context) {
	if !CheckLoginCommandFlags(c) {
		return
	}

	// settings.PrintPlaceholderValues()

	orgInfo := types.OrgInfo{
		AccountName: c.String("org"),
		Region:      c.String("region"),
	}
	password := c.String("password")
	inputUser := c.String("username")
	taURL := ""
	idmServerURL := ""
	userEmail, err := utils.GetUserEmail()
	if err != nil || len(userEmail) == 0 {
		log.Debug(err.Error())
		utils.CheckError(errors.New("Username is not set."))
	}

	if err, value := settings.GetPlaceHolderValue(settings.TIBCO_ACCOUNTS_URL_PLACEHOLDER); err == nil {
		taURL = value
	} else {
		log.Debug(err.Error())
		utils.CheckError(errors.New("TIBCO Accounts URL is not set."))
	}

	idmServerURL, err = utils.GetIDMConnectURL()
	if err != nil {
		log.Debug(err.Error())
		utils.CheckError(errors.New("Identity-Management Server URL is not set."))
	}

	if !c.IsSet("username") && !c.IsSet("password") {
		inputUser = utils.PromptForUser(userEmail)
	} else if !c.IsSet("username") && c.IsSet("password") {
		inputUser = userEmail
	}

	if CheckPlatformVersionAndLogin(c) == nil {
		cOrg, cRegion, err := utils.GetOrgAndRegion()
		if err != nil {
			log.Debug(err.Error())
			utils.CheckError(errors.New("Failed to retrieve current organization and region information. "))
			return
		}

		if inputUser == userEmail {
			if c.IsSet("org") && c.IsSet("region") {
				if cOrg == orgInfo.AccountName && cRegion == orgInfo.Region {
					fmt.Println("User is already logged in. ")
					return
				}
			} else {
				fmt.Println("User is already logged in. ")
				return
			}
		}

	}

	if !c.IsSet("password") {
		password = utils.PromptForPassword()
	}

	// do login
	token, err := TaLogin(taURL, inputUser, password)
	utils.CheckError(err)

	err = IdmLogin(idmServerURL, inputUser, token.AccessToken, orgInfo, true)

	utils.CheckError(err)
	return
}

// IsValidPlatformApi validates cli version against platform api by accessing /platformapiversion
func IsValidPlatformApi() (bool, error) {

	domainUrl, err := utils.GetDomainURL()
	if err != nil {
		return false, err
	}

	parsedDomainURL, err := url.Parse(domainUrl)
	if err == nil {
		parsedDomainURL.Path = "/platformapiversion"
	} else {
		return false, err
	}

	tccUrlForPlatformVersion := parsedDomainURL.String()
	log.Debugf("tccURL to GET PlatformVersion: %s", tccUrlForPlatformVersion)

	noRedirectMarker := errors.New("my-redirect-marker")

	httpClient := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return noRedirectMarker
	}}
	request, err := http.NewRequest("GET", tccUrlForPlatformVersion, nil)
	if err != nil {
		log.Debugf("GET '%s' failed with error: %s", tccUrlForPlatformVersion, err.Error())
		return false, err
	}

	resp, err := httpClient.Do(request)
	if err != nil && !strings.Contains(err.Error(), noRedirectMarker.Error()) {
		log.Debugf("Request '%+v' failed with error: %s", request, err.Error())
		return false, err
	}
	defer resp.Body.Close()

	//there's no version to parse unless we get 200 status code
	if resp.StatusCode != http.StatusOK {
		return false, errors.New("Server responded with " + resp.Status)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debugf("Read resp.Body '%+v' failed with error: %s", resp.Body, err.Error())
		return false, err
	}

	platformApiVersion := string(bytes)

	log.Debugf("PlatformApiVersion received from '%s': %s", tccUrlForPlatformVersion, platformApiVersion)

	splitPlatformApiVersion := strings.Split(platformApiVersion, ".")
	splitCliVersion := strings.Split(consts.CLI_VERSION, ".")

	log.Debugf("platformApiVersion: '%s' consts.CLI_VERSION: '%s'", platformApiVersion, consts.CLI_VERSION)

	// compare the major version of cli and platform api version
	if splitCliVersion[0] != splitPlatformApiVersion[0] {
		log.Debugf("Platform major version '%s' doesn't match CLI major version '%s'", splitPlatformApiVersion[0], splitCliVersion[0])
		return false, nil
	}
	return true, nil
}

// CheckPlatformVersionAndLogin is a before action for all commands which enforces user to be logged-in first
//ensure that platform API version is correct and user is logged in
func CheckPlatformVersionAndLogin(c *cli.Context) error {
	//	validate tibcli api version compatibility with platform api version

	log.Debug("Validate tibcli version and check login...")

	if isValid, err := IsValidPlatformApi(); err != nil {
		if strings.Contains(err.Error(), "Session.orgList") {
			DeleteSessionFile()
			DeleteTokenFile()
			return checkLogin(c)
		}
		utils.CheckError(errors.New(fmt.Sprintf("Troposphere platform services api version validating failed with error : %s", err.Error())))
	} else {
		if isValid == false {
			utils.CheckError(errors.New("Troposphere platform services api version mismatched with tibcli version, please download a new tibcli command line tool from web page."))
		}
	}
	return checkLogin(c)
}

func CheckLoginCommandFlags(c *cli.Context) bool {
	if c.IsSet("username") && !c.IsSet("password") {
		fmt.Println("Please provide password if username is specified. \n ")
		fmt.Println("Example: ")
		fmt.Println("  tib-cli login -u yourname@example.com -p yourpassword \n ")
		return false
	}
	if (c.IsSet("org") && !c.IsSet("region")) || (!c.IsSet("org") && c.IsSet("region")) {
		utils.CheckError(errors.New("Please provide organization name with region info. Either organization name or region is missing. "))
		return false
	}
	return true
}
