// Copyright (c) 2015-2017 TIBCO Software Inc.
// All Rights Reserved

// Package client implements the clients for TIBCO Account's OAuth2 server as well as Domain Server, Build Server, App Manager Web Client
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Morphyni/tas-cli/consts"
	"github.com/Morphyni/tas-cli/settings"
	"github.com/Morphyni/tas-cli/types"
	log "github.com/sirupsen/logrus"
)

// OAuth2 encapsulates the communications with TIBCO Account's OAuth2 server
type OAuth2 interface {
	Login(types.AuthRequest) (*types.OAResponse, error)

	Renew(types.FollowupRequest) (*types.OAResponse, error)

	Logout(types.FollowupRequest) (*types.OAResponse, error)
}

type oAuthClient struct {
	url     *url.URL
	profile *settings.Profile
}

// Enforce implementation of OAuth2 interface
var _ OAuth2 = (*oAuthClient)(nil)

//constructor
func NewOAuthClient(serverURL string) (OAuth2, error) {
	if parsedURL, err := url.Parse(serverURL); err == nil {
		return &oAuthClient{url: parsedURL}, nil
	} else {
		return nil, err
	}
}

//login function
func (client *oAuthClient) Login(rqst types.AuthRequest) (*types.OAResponse, error) {
	data := url.Values{
		"grant_type": {"password"},
		"client_id":  {rqst.ClientId},
		"username":   {rqst.Username},
		"password":   {rqst.Pwd},
		//		"scope"        : {rqst.Scope},
		//		"client_secret": {rqst.ClientSecret},
	}
	return postData(data.Encode(), "application/x-www-form-urlencoded", client.url.String())
}

func (client *oAuthClient) Renew(rqst types.FollowupRequest) (*types.OAResponse, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {rqst.Token},
		//		"scope"        : {rqst.Scope},
		"client_id":     {rqst.ClientId},
		"client_secret": {rqst.ClientSecret},
	}
	return postData(data.Encode(), "application/x-www-form-urlencoded", client.url.String())
}

func (client *oAuthClient) Logout(req types.FollowupRequest) (*types.OAResponse, error) {
	data := fmt.Sprintf("{access_token:\"%v\",client_id:\"%v\",client_secret:\"%v\"}",
		req.Token, req.ClientId, req.ClientSecret)
	return postData(data, "application/json", client.url.String()+"/revoke")
}

func postData(data, contentType, url string) (oaResponse *types.OAResponse, err error) {
	request, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", contentType)

	httpClient := http.DefaultClient
	if strings.HasPrefix(url, "https") {
		err, tciDeploymentName := settings.GetPlaceHolderValue(settings.DEPLOYMENT_NAME_PLACEHOLDER)
		if err != nil {
			log.Debug(err.Error())
			return nil, err
		}
		if len(tciDeploymentName) > 0 {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DisableCompression: true,
				DisableKeepAlives:  true,
			}
			//tr.TLSClientConfig.RootCAs.AppendCertsFromPEM([]byte(cert_DigiCert_SHA2_High_Assurance_Server_CA)) system's are ok
			httpClient = &http.Client{Transport: tr}
			log.Debugf("Using %v to connect to %v ", tr, url)
		} else {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{
					//	RootCAs: x509.NewCertPool() //we're not adding new roots
				},
				DisableCompression: true,
				DisableKeepAlives:  true,
			}
			//tr.TLSClientConfig.RootCAs.AppendCertsFromPEM([]byte(cert_DigiCert_SHA2_High_Assurance_Server_CA)) system's are ok
			httpClient = &http.Client{Transport: tr}
			log.Debugf("Using %v to connect to %v ", tr, url)
		}
	}

	if os.Getenv(consts.TASCLI_DBG) != "" {
		log.Debugf("About to post %v ", request)
	}
	response, err := httpClient.Do(request)

	if err != nil {
		if response != nil {
			request.Close = true
			request.Body.Close()
			response.Body.Close()
		}
		return nil, err
	}
	readBody, _ := ioutil.ReadAll(response.Body)
	request.Body.Close()

	oaResponse = &types.OAResponse{}

	//	time.Sleep(250 * time.Millisecond)//that's a bug in Go but even their own test does that

	if response.StatusCode != 200 {
		log.Debugf("Response unsuccessful: %v %v.", response.Status, string(readBody))
	}
	if readBody != nil {
		if os.Getenv(consts.TASCLI_DBG) != "" {
			log.Debugf("Raw response was: %v", string(readBody))
		}
		if err = json.Unmarshal(readBody, oaResponse); err == nil {
			//we parsed the response but need to check the status code
			if response.StatusCode != 200 {
				log.Debugf("Returning %v", oaResponse.ErrorDesc)
				return nil, errors.New(oaResponse.ErrorDesc)
			}
			//else it's a success; fall through
		} else {
			//can't parse JSON; return HTTP error
			if response.StatusCode != 200 {
				return nil, errors.New(response.Status + " " + string(readBody)) //TODO use error code
			}
			//200 but no JSON ?!
			return nil, errors.New(fmt.Sprintf("Unexpected content: '%+v': %+v", readBody, err))
		}
	} else {
		return nil, errors.New(response.Status)
	}
	if os.Getenv(consts.TASCLI_DBG) != "" {
		log.Debugf("Returinng successfully parsed response: %v", oaResponse)
	}
	return
}

func Debug() {
	log.SetLevel(log.DebugLevel)
}
