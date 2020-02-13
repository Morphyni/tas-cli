package consts

const (
	CLI_MODULE_NAME       string = "tas-cli"
	CLI_MODULE_USAGE      string = "TIBCO AuditSafe Cloud - Command Line Interface"
	CLI_VERSION           string = "1.0.0"
	CLI_COPYRIGHT_MESSAGE string = "2015-2019 TIBCO Software Inc."

	//environment property for indicating extra debug info; sensitive data will be displayed. Don't use
	TASCLI_DBG string = "TASCLI_DBG"

	OBFUSCATE_COOKIE_VALUE = true
	// Context paths of all WebClient's
	WEB_SERVER_CONTEXT_PATH          string = "/api/"
	DOMAIN_SERVER_CONTEXT_PATH       string = "/domain/"
	APP_MANAGER_CONTEXT_PATH         string = "/appmgr/"
	ORCHESTRATOR_CONTEXT_PATH        string = "/orchestrator/"
	APPLOG_CONTEXT_PATH              string = "/tropos/unityapiproxy/"
	IDENTITY_MANAGEMENT_CONTEXT_PATH string = "/idm/"
	BUILD_SERVER_CONTEXT_PATH        string = "/build/"

	// API versions of all WebClient's
	WEB_SERVER_API_VERSION          string = "v1"
	DOMAIN_SERVER_API_VERSION       string = "v1"
	DOMAIN_SERVER_API_VERSION_V2    string = "v2"
	APP_MANAGER_API_VERSION         string = "v1"
	ORCHESTRATOR_API_VERSION        string = "v1"
	APPLOG_VERSION                  string = "v1"
	IDENTITY_MANAGEMENT_API_VERSION string = "v2"
	BUILD_SERVER_API_VERSION        string = "v1"
	// IDENTITY_MANAGEMENT_API_VERSION string = "v1"

	// Orchestrator WebClient API's
	ORCHESTRATOR_SANDBOXES_API     string = "/sandboxes"
	ORCHESTRATOR_APPS_API          string = "/apps"
	ORCHESTRATOR_APP_API           string = "/app"
	ORCHESTRATOR_STATUS_API        string = "/status"
	ORCHESTRATOR_INSTANCES_API     string = "/instances"
	ORCHESTRATOR_CONFIGURATION_API string = "/configuration"
	ORCHESTRATOR_FTL_STATUS        string = "/ftl"
	ORCHESTRATOR_APPLICATIONS_API  string = "/applications"

	// Domain Server WebClient API's
	DOMAIN_SERVER_LOGIN_API string = "/login-oauth"

	// IdentityManagementServer WebClient API's
	IDENTITY_MANAGEMENT_LOGIN_API string = "/login-oauth"

	DOMAIN_SERVER_SANDBOXES_API string = "/sandboxes"

	DOMAIN_SERVER_USERS_API string = "/users"

	// App Manager WebClient API's
	APP_MANAGER_APPS_API string = "/apps"

	// Web Server API's
	WEB_SERVER_FEATURE_API string = "/feature"

	// Atmosphere Cookie name
	//ATMOS_COOKIE_NAME string = "AtmosphereSession"

	// Common statuses
	SUCCESS_STATUS   string = "success"
	ERROR_STATUS     string = "error"
	COMPLETED_STATUS string = "completed"
	//environment property governing persistence of OAuth tokens
	DONT_PERSIST string = "TIBCLI_DONT_PERSIST"

	//Hostname and port settings for the current local envrioment.
	WEBAPI_LOCAL_HOST = "http://localhost"
	WEBAPI_LOCAL_PORT = "3000"

	APPLOG_QUERY_API string = "/logsQueries"
)

const (
	DEFAULT_SANDBOX = "MyDefaultSandbox"
)
