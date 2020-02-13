package commands

import (
	"errors"
	"fmt"

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
