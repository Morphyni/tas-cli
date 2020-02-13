// Copyright (c) 2015-2017 TIBCO Software Inc.
// All Rights Reserved

// package utils contains common utilities
package utils

import (
	"net/url"

	"github.com/Morphyni/tas-cli/consts"
)

// GetDomainServerGetSandboxesAPI returns REST API path for Domain-Server to get all sandboxes including operational
func GetDomainServerGetOpSandboxesAPI() string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + consts.DOMAIN_SERVER_SANDBOXES_API
}

// GetWebServerFeatureAPI returns REST API path for Domain-Server login
func GetWebServerFeatureAPI() string {
	return consts.WEB_SERVER_CONTEXT_PATH + consts.WEB_SERVER_API_VERSION + consts.WEB_SERVER_FEATURE_API
}

// GetDomainServerGetAppAPI returns REST API path for Domain-Server get app
func GetDomainServerGetAppAPI(sandboxId string) string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + consts.DOMAIN_SERVER_SANDBOXES_API + "/" + sandboxId + "/applications"
}

// GetDomainServerGetAppDetailsAPI returns REST API path for Domain-Server get app details
func GetDomainServerGetAppDetailsAPI(sandboxId string, appId string) string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + consts.DOMAIN_SERVER_SANDBOXES_API + "/" + sandboxId + "/applications/" + appId
}

// GetDomainServerGetAppConfigAPI returns REST API path for Domain-Server get app configuration
func GetDomainServerGetAppConfigAPI(sandboxId, appId string) string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + consts.DOMAIN_SERVER_SANDBOXES_API + "/" + sandboxId + "/applications/" + appId + "/configuration"
}

// GetDomainServerListAppsAPI returns REST API path for Domain-Server List Apps
func GetDomainServerListAppsAPI(sandboxId string) string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + consts.DOMAIN_SERVER_SANDBOXES_API + "/" + sandboxId + "/applications"
}

// GetDomainServerListAllAppsAPI returns REST API path for Domain-Server List All Apps
func GetDomainServerListAllAppsAPI() string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + "/applications"
}

// GetDomainServerGetSandboxesAPI returns REST API path for Domain-Server get all sandboxes
func GetDomainServerGetSandboxesAPI() string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + consts.DOMAIN_SERVER_SANDBOXES_API
}

// GetDomainServerDefaultSandboxAPI returns REST API path for Domain-Server get all sandboxes
func GetDomainServerDefaultSandboxAPI() string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + consts.DOMAIN_SERVER_SANDBOXES_API + "/" + consts.DEFAULT_SANDBOX
}

// GetDomainServerGetSandboxAPI returns REST API path for Domain-Server get a sandbox
func GetDomainServerGetSandboxAPI(sandboxId string) string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + consts.DOMAIN_SERVER_SANDBOXES_API + "/" + sandboxId
}

// GetDomainServerGetAppEndpoint returns REST API path for Domain-Server get endpoint
func GetDomainServerGetAppEndpointAPI(sandboxId string, applicationId string, endpointId string) string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + consts.DOMAIN_SERVER_SANDBOXES_API + "/" + sandboxId + "/applications/" + applicationId + "/endpoints/" + endpointId
}

// GetDomainServerGetAppEndpointURL returns REST API path for Domain-Server get endpoint URL
func GetDomainServerGetAppEndpointURLAPI(sandboxId string, applicationId string, endpointId string) string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + consts.DOMAIN_SERVER_SANDBOXES_API + "/" + sandboxId + "/applications/" + applicationId + "/endpoints/" + endpointId + "/url"
}

// GetDomainServerFetchAppAuditsAPI returns REST API path for Domain-Server fetch app audit history URL
func GetDomainServerFetchAppAuditsAPI(applicationId string) string {
	return consts.DOMAIN_SERVER_CONTEXT_PATH + consts.DOMAIN_SERVER_API_VERSION + "/audits/" + applicationId
}

// GetAppManagerNewAppsAPI returns REST API path for App-Manager to start new application
func GetAppManagerNewAppsAPI(appId string) string {
	return consts.APP_MANAGER_CONTEXT_PATH + consts.APP_MANAGER_API_VERSION + consts.APP_MANAGER_APPS_API + "/" + appId
}

// GetOrchestratorSandboxAPI returns REST API path for Orchestrator to push application
func GetOrchestratorSandboxAPI() string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_SANDBOXES_API
}

// GetOrchestratorPushAPI returns REST API path for Orchestrator to push application
func GetOrchestratorPushAPI(sandboxId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_SANDBOXES_API + "/" + sandboxId + consts.ORCHESTRATOR_APPS_API
}

// GetOrchestratorMoveAppAPI returns REST API path for Orchestrator to move application from one sandbox to another
func GetOrchestratorMoveAppAPI(sandboxId, appId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_SANDBOXES_API + "/" + sandboxId + "/applications/" + appId + "/target"
}

// GetOrchestratorPromoteAppAPI returns REST API path for Orchestrator to promote application to operational sandbox
func GetOrchestratorPromoteAppAPI(sandboxId, appId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_SANDBOXES_API + "/" + sandboxId + "/applications/" + appId + "/promote"
}

// GetOrchestratorCopyAppAPI copies app
func GetOrchestratorCopyAppAPI(appId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + "/applications/" + appId + "/copy"
}

// GetOrchestratorUpgradeAppAPI returns REST API path for Orchestrator to upgrade application to with another application from operational sandbox
func GetOrchestratorUpgradeAppAPI() string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_APP_API + "/upgrade"
}

// GetOrchestratorReplaceAppAPI returns REST API path for Orchestrator to replace application with another application
func GetOrchestratorReplaceAppAPI(appId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + "/applications/" + appId + "/replace"
}

// GetOrchestratorUpdateAppAPI returns REST API path for Orchestrator to update application attributes with given values
func GetOrchestratorUpdateAppAPI(sandboxId, appId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_SANDBOXES_API + "/" + sandboxId + "/applications/" + appId
}

// GetOrchestratorUpdateApplicationAPI returns REST API path for Orchestrator to update application visibility
func GetOrchestratorUpdateAppVisibilityAPI(appId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_APPLICATIONS_API + "/" + appId
}

// GetOrchestratorConfigureAPI returns REST API path for Orchestrator to configure application property overrides
func GetOrchestratorConfigureAPI(sandboxId, appId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_SANDBOXES_API + "/" + sandboxId + consts.ORCHESTRATOR_APPS_API + "/" + appId + consts.ORCHESTRATOR_CONFIGURATION_API
}

// GetOrchestratorDeleteAPI returns REST API path for Orchestrator to delete application
func GetOrchestratorDeleteAPI(sandboxId, appId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_SANDBOXES_API + "/" + sandboxId + consts.ORCHESTRATOR_APPS_API + "/" + appId
}

// GetOrchestratorStatusAPI returns REST API path for Orchestrator to get application status
func GetOrchestratorStatusAPI(sandboxId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_SANDBOXES_API + "/" + sandboxId + consts.ORCHESTRATOR_APPS_API + consts.ORCHESTRATOR_STATUS_API
}

// GetOrchestratorScaleAPI returns REST API path for Orchestrator to scale application
func GetOrchestratorScaleAPI(sandboxId, appId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_SANDBOXES_API + "/" + sandboxId + consts.ORCHESTRATOR_APPS_API + "/" + appId + consts.ORCHESTRATOR_INSTANCES_API
}

// GetBuildServerPushBWSupplementAPI return API path for BuildServer to push supplement
func GetBuildServerPushBWSupplementAPI() string {
	return consts.BUILD_SERVER_CONTEXT_PATH + consts.BUILD_SERVER_API_VERSION + "/supplement"
}

// GetAppLogCreateQueryAPI returns REST API path for AppLog to create a query.
func GetAppLogCreateQueryAPI() string {
	return consts.APPLOG_CONTEXT_PATH + consts.APPLOG_VERSION + consts.APPLOG_QUERY_API
}

// GetAppLogGetLogsAPI returns REST API path for AppLog to fetch application logs.
func GetAppLogGetLogsAPI(queryId string) string {
	return consts.APPLOG_CONTEXT_PATH + consts.APPLOG_VERSION + consts.APPLOG_QUERY_API + "/" + queryId
}

// GetAppLogDeleteQueryAPI returns REST API path for AppLog to delete a log query.
func GetAppLogDeleteQueryAPI(queryId string) string {
	return consts.APPLOG_CONTEXT_PATH + consts.APPLOG_VERSION + consts.APPLOG_QUERY_API + "/" + queryId
}

// GetIdentityManagementLoginAPI returns REST API path for Identity-Management login
func GetIdentityManagementLoginAPI() string {
	return consts.IDENTITY_MANAGEMENT_CONTEXT_PATH + consts.IDENTITY_MANAGEMENT_API_VERSION + consts.IDENTITY_MANAGEMENT_LOGIN_API
}

// GetFTLStatus method gets FTL Status of a sandbox from Atmosphere(Orchestrator)
func GetFTLStatus(sandboxId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + consts.ORCHESTRATOR_FTL_STATUS + consts.DOMAIN_SERVER_SANDBOXES_API + "/" + sandboxId
}

// GetOrchestratorEnableDisableOrgFTLAPI returns REST API path for Orchestrator to enable/disable FTL
func GetEnableDisableOrgFTLAPI() string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + "/organization/ftl"
}

// GetOrgInfoURLAPI returns tci apps domain
func GetOrgInfoURLAPI() string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + "/orginfo"
}

// GetAcessKeysAPI returns tci apps domain
func GetAcessKeysAPI() string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + "/tibtunnel/accesskeys"
}

// GetUpdateAppAccessKeyAPI returns API path for update accessKey for app
func GetUpdateAppAccessKeyAPI(sandboxId, appId, accessKey string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + "/sandboxes/" + sandboxId + "/applications/" + appId + "/tibTunnelAccessKey/" + accessKey
}

// GetRemoveAppAccessKeyAPI returns API path for remove accessKey from app
func GetRemoveAppAccessKeyAPI(sandboxId, appId string) string {
	return consts.ORCHESTRATOR_CONTEXT_PATH + consts.ORCHESTRATOR_API_VERSION + "/sandboxes/" + sandboxId + "/applications/" + appId + "/tibTunnelAccessKey"
}

// ResolveAppLogServerURL returns the base URL for AppLogServer REST API
func ResolveAppLogServerURL() (string, error) {
	if GetEnvParam("WEBAPI_TEST_MODE") == "local" {
		return consts.WEBAPI_LOCAL_HOST + ":" + consts.WEBAPI_LOCAL_PORT, nil
	}
	profile, err := LoadProfile()
	if err != nil {
		return "", err
	}
	u, err := url.Parse(profile.IDMConnectURL)
	if err != nil {
		return "", err
	}
	return u.Scheme + "://" + u.Host, nil
}
