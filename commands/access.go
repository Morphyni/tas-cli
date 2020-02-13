package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Morphyni/tas-cli/client"
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

// TaLogin performs login to TIBCO Accounts with username and password
func TaLogin(url, user, password string) (*types.OAResponse, error) {

	oauth, err := client.NewOAuthClient(url)
	if err != nil {
		return nil, err
	}
	err, clientId := settings.GetPlaceHolderValue(settings.TIBCO_ACCOUNTS_CLIENTID_PLACEHOLDER)
	if err != nil {
		log.Debug(err)
		utils.CheckError(errors.New("TIBCO Accounts' client id not set"))
	}

	resp, err := oauth.Login(types.AuthRequest{Username: user, Pwd: password, ClientId: clientId})
	if err == nil {
		token, e := settings.NewToken()
		if e != nil {
			log.Errorf("NON-FATAL: Couldn't create session file for login token: %v", e)
		}

		// here we keep TA accessToken in a Cookie just want to get benefit of reusing the isValid() func in settingsfile.go which checks the cookie expired or not.
		taTokenCookie := &http.Cookie{Name: settings.ACCESS_TOKEN_KEY_NAME, Value: resp.AccessToken, //TODO switch to RefreshToken
			Expires: time.Now().UTC().Add(time.Duration(resp.ExpiresIn) * time.Second)}

		token.AccessToken = taTokenCookie

		if utils.GetEnvParam(consts.DONT_PERSIST) == "" {
			e = token.Write(consts.OBFUSCATE_COOKIE_VALUE) //obfuscate Token
			if e != nil {
				log.Errorf("NON-FATAL: Couldn't persist the login token to disk: %v", e)
			} else {
				log.Debugf("Persisted the login token to disk.")
			}
		} else {
			log.Infof("No OAuth token persisted since environment variable '%s' wasn't set.", consts.DONT_PERSIST)
		}
	}
	return resp, err
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

//ensure user is logged in
func checkLogin(c *cli.Context) error {

	//check that branded values in tibcli, including domainURL match saved local profile
	profile, session, token, err := utils.LoadSettings()
	utils.CheckError(err)

	// if 'old-style' profile doesn't have version field or version field has wrong version number then wipe out all profile/session/token and start from scratch
	if profile.Version != consts.CLI_VERSION {

		utils.CheckError(profile.Delete())
		utils.CheckError(session.Delete())
		utils.CheckError(token.Delete())

		session, err = utils.LoadSession(consts.OBFUSCATE_COOKIE_VALUE)
		utils.CheckError(err)
		profile, err = utils.LoadProfile()
		utils.CheckError(err)
		token, err = utils.LoadToken(consts.OBFUSCATE_COOKIE_VALUE)
		utils.CheckError(err)
	}

	// get idm server url from placeholder
	err, idmServerURL := settings.GetPlaceHolderValue(settings.IDENTITY_MANAGEMENT_SERVER_HOST_PLACEHOLDER)
	if err != nil || len(idmServerURL) == 0 {
		log.Debug(err.Error())
		utils.CheckError(errors.New("Identity-Management Server URL is not set via placeholder."))
	}

	// get user email from placeholder
	err, userEmail := settings.GetPlaceHolderValue(settings.USERNAME_PLACEHOLDER)
	if err != nil || len(userEmail) == 0 {
		log.Debug(err.Error())
		utils.CheckError(errors.New("Username is not set via placeholder."))
	}

	// check user email in profile
	if len(profile.UserEmail) > 0 && profile.UserEmail != userEmail {
		log.Debugf("profile.userEmail: %s, userEmail read from placeholder is: %s", profile.UserEmail, userEmail)
		userEmail = profile.UserEmail
	}

	// check idm server url in profile
	if len(profile.IDMConnectURL) > 0 && profile.IDMConnectURL != idmServerURL {
		log.Debugf("profile.idmServerUrl: %s, idm server URL read from placeholder is: %s", profile.IDMConnectURL, idmServerURL)
		idmServerURL = profile.IDMConnectURL
	}

	//we may prompt user here for password only if this is NOT the login command.  That command does prompt on its own
	promptForPassword := c.Command.Name != "login"

	isSessionValid := CheckCookiesIsValid(session.Cookies, settings.SESSION_FILENAME)
	log.Debugf("session cookies is valid ? : '%+v'", isSessionValid)

	// check cookies of session are still-valid or not
	if !isSessionValid {
		log.Debug("session cookies get expired")
		// load TA access token
		loadedAccessToken, er := utils.LoadToken(consts.OBFUSCATE_COOKIE_VALUE)
		utils.CheckError(er)

		if CheckCookiesIsValid([]*http.Cookie{loadedAccessToken.AccessToken}, settings.TOKEN_FILE_NAME) {
			log.Debug("Session cookies missing or expired, trying still-valid access token")

			// Login to IDM again to refresh session with the still-valid TA access token
			// Re-login, user won't input username/password, so LoginFlag set to false
			orgInfo := types.OrgInfo{
				AccountName: c.String("org"),
				Region:      c.String("region"),
			}
			err = IdmLogin(idmServerURL, userEmail, loadedAccessToken.AccessToken.Value, orgInfo, false)
			if err != nil {
				log.Debugf("Refresh session rejected by IDM Server %+v ", err)
				utils.CheckError(err)
				return err
			}

		} else { // TA access token and IDM session both expired
			if !promptForPassword { // for login command
				return errors.New("No session cookies, no access token; need to log in")
			}

			log.Debug("User is not logged-in or session has expired. Hence initiating login.")
			fmt.Println("User is not logged-in or session has expired.")

			// get TA URL from placeholder
			err, taURL := settings.GetPlaceHolderValue(settings.TIBCO_ACCOUNTS_URL_PLACEHOLDER)
			if err != nil {
				log.Debug(err.Error())
				utils.CheckError(errors.New("TIBCO Accounts url not set via placeholder."))
			}
			// read user/password from user input
			userEmail = utils.PromptForUser(userEmail)
			password := utils.PromptForPassword()

			// Do TA login
			accessToken, err := TaLogin(taURL, userEmail, password)
			if err != nil {
				log.Debugf("Refresh accessToken rejected by TA Server %+v ", err)
				utils.CheckError(err)
				return err
			}
			// Do IDM login
			err = IdmLogin(idmServerURL, userEmail, accessToken.AccessToken, types.OrgInfo{}, true)
			if err != nil {
				log.Debugf("Refresh session rejected by IDM Server %+v ", err)
				utils.CheckError(err)
				return err
			}
		}
	}
	return nil
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

// CheckCookiesIsValid checks all cookies are valid/expired
func CheckCookiesIsValid(cookies []*http.Cookie, fileName string) bool {
	if cookies == nil || len(cookies) == 0 {
		log.Debug("Cookies of session file are empty")
		return false
	}
	if fileName == settings.TOKEN_FILE_NAME {
		for _, cookie := range cookies {
			if cookie == nil {
				log.Debug("One of the session cookies is empty")
				return false
			}
			if !cookie.Expires.IsZero() && cookie.Expires.Before(time.Now().UTC()) {
				// Cookie expiration is set and before current time means expired
				return false
			}
		}
	} else if fileName == settings.SESSION_FILENAME { // cookies in session file are coming from IDM which don't have 'Expires' field so we need another way to test they are valid or not
		// we use the get default sandbox call on domain server to check if the cookies
		dsClient, err := client.NewDomainServer()
		if err != nil {
			log.Errorf("Initializing DomainServer client instance on error: %s", err.Error())
			return false
		}
		_, httpCode, err := dsClient.GetDefaultSandbox()
		if err != nil {
			if httpCode == 599 && strings.Contains(err.Error(), "419") { // the backend code set 599 code in the response httpCode and the 419 code can be found from error msg
				log.Debugf("Cookies in session file get invalid as we got 419 error while accessing backend: %s", err.Error())
			} else {
				log.Debugf("Validating session cookies failed on errors other then 419: %s", err.Error())
			}
			return false
		}
	} else {
		log.Errorf("Unknown file name: %s", fileName)
		return false
	}
	return true
}

// DeleteSessionFile delete session file from disk
func DeleteSessionFile() {
	session, err := settings.NewSession()
	utils.CheckError(err)
	err = session.Delete()
	utils.CheckError(err)
}

// DeleteTokenFile delete token file from disk
func DeleteTokenFile() {
	token, err := settings.NewToken()
	utils.CheckError(err)
	err = token.Delete()
	utils.CheckError(err)
}

// DeleteProfile delete profile file from disk
func DeleteProfile() {
	profile, err := settings.NewProfile()
	utils.CheckError(err)
	err = profile.Delete()
	utils.CheckError(err)
}
