// Copyright (c) 2015-2017 TIBCO Software Inc.
// All Rights Reserved

// Package client implements the clients for Domain Server, Build Server, App Manager Web Client
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Morphyni/tas-cli/types"
	"github.com/Morphyni/tas-cli/utils"
	log "github.com/sirupsen/logrus"
)

// DomainServer encapsulates the remote operations with the Atmosphere domain server
type DomainServer interface {
	// UpdateSandbox method updates a sandbox in the Domain Server
	// UpdateSandbox(sandboxId string, bodyComponents []byte) (*types.SuccessResponse, error)

	// //GetOrgSandboxes method retrieves sandboxes of current organization
	// GetOrgSandboxes() (*types.DomainServerGetSandboxesResponse, error)

	// //GetDefaultSandbox method provides default sandbox of current organization
	GetDefaultSandbox() (*types.DomainServerSandboxBean, int, error)

	// // GetSandboxes method retrieves all sandboxes for the current user from Domain Server
	// GetSandboxes(sandboxName, userName string) (*types.DomainServerGetSandboxesResponse, error)

	// // GetSandbox method retrieves a sandbox from Domain Server
	// GetSandbox(sandboxId string) (*types.DomainServerSandboxBean, error)

	// // GetApplicationsInSandbox method retrieves the Applications Beans from Domain Server
	// GetApplicationsInSandbox(sandboxId string) (*types.DomainServerApplicationsResponse, error, bool)

	// // GetAllApplications method retrieves the all Applications Beans from Domain Server
	// GetAllApplications() (*types.DomainServerApplicationsResponse, error, bool)

	// // GetApplicationDetails method retrieves an App from Domain Server, the last boolean return argument is true if the error is a simple application not found on that sandbox
	// GetApplicationDetails(appName, sandboxId string) (*types.DomainServerApplicationBean, error, bool)

	// // GetAppEndpoint method retrieves an App endpoint bean from Domain Server
	// GetAppEndpoint(sandboxId, appId, endpointId string) (*types.DomainServerAppEndpointBean, error)

	// // GetAppEndpointUrl method retrieves an App endpoint URL from Domain Server
	// GetAppEndpointUrl(sandboxId, appId, endpointId string) (*types.DomainServerAppEndpointUrlResponse, error)

	// // GetAppConfigDetails gets app config details
	// GetAppConfigDetails(sandboxId, appId string) (*types.AppConfig, error)

	// // GetAllApplicationsBySbscId retrieves all apps belong to targetSbsc
	// GetAllApplicationsBySbscId(targetSbscId string) (*types.DomainServerApplicationsResponse, error)

	// // GetApps return apps by appName and jwt:role
	// GetApps(appName string, role bool) ([]types.DomainServerApplicationBean, error)

	// // GetApp returns app by appId
	// GetApp(appId string) (*types.DomainServerApplicationBean, error)

	// // GetAppAudits
	// GetAppAudits(appId string) (*types.DomainServerAppAudits, error)
}

// domainServer is the private implementation of the DomainServer interface
type domainServer struct {
	webClient
}

// make sure that the domainServer implements the DomainServer interface
var _ DomainServer = (*domainServer)(nil)

// NewDomainServer creates a new DomainServer object
func NewDomainServer() (ds DomainServer, err error) {
	serverURL, err := utils.GetDomainURL()
	if err != nil {
		return nil, err
	}
	if parsedURL, err := url.Parse(serverURL); err == nil {
		return &domainServer{webClient: webClient{url: parsedURL}}, nil
	} else {
		return nil, err
	}
}

// NewDomainServerV2 creates a new DomainServer object
func NewDomainServerV2(domainURL string) (ds DomainServer, err error) {
	var serverURL string
	if len(domainURL) > 0 {
		serverURL = domainURL
	} else {
		serverURL, err = utils.GetDomainURL()
		if err != nil {
			return nil, err
		}
	}
	if parsedURL, err := url.Parse(serverURL); err == nil {
		return &domainServer{webClient: webClient{url: parsedURL}}, nil
	}
	return nil, err
}

func (c *domainServer) GetDefaultSandbox() (*types.DomainServerSandboxBean, int, error) {

	c.url.Path = utils.GetDomainServerDefaultSandboxAPI()

	log.Debugf("Sending GetDefaultSandbox against url: '%s'", c.url.String())

	response, restCallAndCookieRefreshErr := utils.RestCallAndCookiesRefreshHandler(
		&types.RestCallRequest{
			Url:          c.url,
			Headers:      map[string]string{"Content-Type": "application/json"},
			Method:       http.MethodGet,
			Body:         nil,
			LogRequest:   false,
			UserId:       "",
			RetryAttempt: nil,
		})

	if restCallAndCookieRefreshErr != nil {
		log.Debugf("Rest call and refresh cookies failed on error: '%+v'", restCallAndCookieRefreshErr.Error())
		return nil, http.StatusInternalServerError, restCallAndCookieRefreshErr
	}

	if utils.IsDevMode() {
		log.Debugf("Response for DomainServer GetDefaultSandbox request: %+v", response)
	}

	if response.ErrorResponse != nil {
		atmosError := response.ErrorResponse
		errMsg := fmt.Sprintf("Error for GET '%s' url: ErrorCode: %s, ErrorMsg: %s, ErrorDetail: %s .",
			c.url.String(), atmosError.ErrorCode, atmosError.ErrorMsg, atmosError.ErrorDetail)
		log.Debug(errMsg)
		return nil, response.HttpCode, errors.New(atmosError.ErrorMsg)
	}

	sandboxBean := &types.DomainServerSandboxBean{}
	//Populate a structure type with the byte data returned by the REST Call
	err3 := json.Unmarshal(response.ResponseBytes, &sandboxBean)
	if err3 != nil {
		if utils.IsDevMode() {
			utilities.PrintStackTrace(2, err3.Error())
		}
		return nil, http.StatusInternalServerError, err3
	}

	return sandboxBean, response.HttpCode, nil
}
