// Copyright (c) 2015-2017 TIBCO Software Inc.
// All Rights Reserved

package settings

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Morphyni/tas-cli/consts"
	log "github.com/sirupsen/logrus"
)

//Placeholder type
const (
	USERNAME_PLACEHOLDER = iota
	DOMAIN_SERVER_HOST_PLACEHOLDER
	TIBCO_ACCOUNTS_URL_PLACEHOLDER
	TIBCO_ACCOUNTS_CLIENTID_PLACEHOLDER
	IDENTITY_MANAGEMENT_SERVER_HOST_PLACEHOLDER
	TCI_TENANT_ID_PLACEHOLDER
	REGION_PLACEHOLDER
	DEPLOYMENT_NAME_PLACEHOLDER
	FTL_ENABLED_OPTION_PLACEHOLDER
)

var PLACEHOLDER_NAMES = [...]string{
	"username_placeholder",
	"domain_server_host_port_placeholder",
	"TA_server_url_placeholder",
	"oauth_cid",
	"identity_management_server_host_port_placeholder",
	"tci_tenant_id_placeholder",
	"region_placeholder",
	"deployment_name_placeholder",
	"ftl_enabled_option_placeholder",
}
var PLACEHOLDER_CHARS = [...]byte{'@', '!', '/', '&', '*', '%', '?', '$', '|'}

var PLACEHOLDER_VALUES = [...]string{
	"@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@",
	"!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!",
	"//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////",
	"&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&",
	"*******************************************************************************************************************",
	"%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%",
	"???????????????????????????????????????????????????????????????????????????????????????????????????????????????????",
	"$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$",
	"|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||"}

//GetPlaceHolderValue returns placeholder value corresponding to the given placeholder type
func GetPlaceHolderValue(placeholderType int) (error, string) {

	// switch placeholderType {
	// case USERNAME_PLACEHOLDER:
	// 	return nil, "jack@testusers.tci.com"
	// case DOMAIN_SERVER_HOST_PLACEHOLDER:
	// 	return nil, "https://d1-localdev-integration.cloud.tibco.com/"
	// case IDENTITY_MANAGEMENT_SERVER_HOST_PLACEHOLDER:
	// 	return nil, "https://d1-localdev-integration.cloud.tibco.com/"
	// case TIBCO_ACCOUNTS_URL_PLACEHOLDER:
	// 	return nil, "https://d1-localdev-integration.cloud.tibco.com/testusers/oauth/token"
	// case TIBCO_ACCOUNTS_CLIENTID_PLACEHOLDER:
	// 	return nil, "ropc_ipass"
	// case TCI_TENANT_ID_PLACEHOLDER:
	// 	return nil, "tci"
	// case REGION_PLACEHOLDER:
	// 	return nil, "vagrant"
	// case DEPLOYMENT_NAME_PLACEHOLDER:
	// 	return nil, ""
	// case FTL_ENABLED_OPTION_PLACEHOLDER:
	//  return nil, "true"
	// }

	if placeholderType != TIBCO_ACCOUNTS_CLIENTID_PLACEHOLDER || os.Getenv(consts.TASCLI_DBG) != "" {
		//don't debug print out sensitive TIBCO_ACCOUNTS_CLIENTID_PLACEHOLDER unless explicitly told so
		//log.Debugf("Placeholder '%s' length '%d' , value: '%s'", PLACEHOLDER_NAMES[placeholderType], len(PLACEHOLDER_VALUES[placeholderType]), PLACEHOLDER_VALUES[placeholderType])
	}

	trimmedValue := strings.TrimSpace(PLACEHOLDER_VALUES[placeholderType])
	if placeholderType != DEPLOYMENT_NAME_PLACEHOLDER && strings.EqualFold(trimmedValue, "") {
		return errors.New(fmt.Sprintf("The place holder '%s' is empty.", PLACEHOLDER_NAMES[placeholderType])), ""
	} else {
		if isPlaceholderReplaced(PLACEHOLDER_CHARS[placeholderType], PLACEHOLDER_VALUES[placeholderType]) {
			//log.Debugf("Returned value of placeholder '%s' is: '%s'", PLACEHOLDER_NAMES[placeholderType], trimmedValue)
			return nil, trimmedValue
		} else {
			return errors.New(fmt.Sprintf("The place holder '%s' is not replaced.", PLACEHOLDER_NAMES[placeholderType])), PLACEHOLDER_VALUES[placeholderType]
		}
	}
}

//isPlaceholderReplaced validates if the placeholder replaced or not
func isPlaceholderReplaced(ch byte, value string) bool {
	for _, v := range []byte(value) {
		if byte(v) != ch {
			return true
		}
	}
	//For the case like 254 '@' chars replaced with 5 '@' will be considered as 'not replaced', false returned
	return false
}

//PrintPlaceholderValues show up the place holder values in log dug which could be used by QA to inspect the values while testing
func PrintPlaceholderValues() {

	_, value := GetPlaceHolderValue(USERNAME_PLACEHOLDER)
	log.Debugf("Placeholder '%s' value: '%s'", PLACEHOLDER_NAMES[USERNAME_PLACEHOLDER], value)

	_, value = GetPlaceHolderValue(DOMAIN_SERVER_HOST_PLACEHOLDER)
	log.Debugf("Placeholder '%s' value: '%s'", PLACEHOLDER_NAMES[DOMAIN_SERVER_HOST_PLACEHOLDER], value)

	_, value = GetPlaceHolderValue(TIBCO_ACCOUNTS_URL_PLACEHOLDER)
	log.Debugf("Placeholder '%s' value: '%s'", PLACEHOLDER_NAMES[TIBCO_ACCOUNTS_URL_PLACEHOLDER], value)

	_, value = GetPlaceHolderValue(TIBCO_ACCOUNTS_CLIENTID_PLACEHOLDER)
	if os.Getenv(consts.TASCLI_DBG) != "" {
		log.Debugf("Placeholder '%s' value: '%s'", PLACEHOLDER_NAMES[TIBCO_ACCOUNTS_CLIENTID_PLACEHOLDER], value)
	}

	_, value = GetPlaceHolderValue(IDENTITY_MANAGEMENT_SERVER_HOST_PLACEHOLDER)
	log.Debugf("Placeholder '%s' value: '%s'", PLACEHOLDER_NAMES[IDENTITY_MANAGEMENT_SERVER_HOST_PLACEHOLDER], value)

	_, value = GetPlaceHolderValue(TCI_TENANT_ID_PLACEHOLDER)
	log.Debugf("Placeholder '%s' value: '%s'", PLACEHOLDER_NAMES[TCI_TENANT_ID_PLACEHOLDER], value)

	_, value = GetPlaceHolderValue(REGION_PLACEHOLDER)
	log.Debugf("Placeholder '%s' value: '%s'", PLACEHOLDER_NAMES[REGION_PLACEHOLDER], value)

	_, value = GetPlaceHolderValue(DEPLOYMENT_NAME_PLACEHOLDER)
	log.Debugf("Placeholder '%s' value: '%s'", PLACEHOLDER_NAMES[DEPLOYMENT_NAME_PLACEHOLDER], value)

	_, value = GetPlaceHolderValue(FTL_ENABLED_OPTION_PLACEHOLDER)
	log.Debugf("Placeholder '%s' value: '%s'", PLACEHOLDER_NAMES[FTL_ENABLED_OPTION_PLACEHOLDER], value)
}
