package types

import (
	"bytes"
	"io"
	"net/url"

	"git.tibco.com/git/product/troposphere/golib.git/atmos/utilities"
	"github.com/urfave/cli"
)

//source and target sandbox details
type UpgradeAppInfo struct {
	SourceSandboxId string `json:"sourceSandboxId"`
	SourceAppId     string `json:"sourceAppId"`
	TargetSandboxId string `json:"targetSandboxId"`
	TargetAppId     string `json:"targetAppId"`
}

//loggedInUserInfo
type LoggedInUserInfo struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

//accountsInfo
type AccountsInfo struct {
	AccountId          string          `json:"accountId"`
	AccountDisplayName string          `json:"accountDisplayName"`
	LoggedInUserRole   string          `json:"loggedInUserRole"`
	RegionToUrls       []RegionUrlInfo `json:"regionsToUrls"`
	SubscriptionId     string          `json:"subscriptionId"`
}

//regionsToUrls
type RegionUrlInfo struct {
	Region string `json:"region"`
	Url    string `json:"url"`
}

// MultiSubscriptionLoginResponse is response received from domain server for Multi Subscription login request
type MultiSubscriptionLoginResponse struct {
	Accounts     []AccountsInfo   `json:"accountsInfo"`
	LoggedInUser LoggedInUserInfo `json:"loggedInUserInfo"`
}

//OrgInfo to store account name for multi subscription org details
type OrgDetails struct {
	AccountName    string
	Region         string
	AccountId      string
	RegionUrl      string
	SubscriptionId string
}

//OrgInfo to store account name for multi subscription login request
type OrgInfo struct {
	AccountName string
	Region      string
}

// CommandLineContext is an interface for the cli.Context to be able to mock objects
type CommandLineContext interface {
	Args() cli.Args
	String(string) string
	IsSet(string) bool
	Bool(string) bool
}

// IDMLoginResponse is response received from domain server for login request
type IDMLoginResponse struct {
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	UserName       string     `json:"userName"`
	UserId         string     `json:"userId"`
	OrgName        string     `json:"orgName"`
	TS             int        `json:"ts"`
	DomainUrl      string     `json:"domainUrl"`
	OrgDisplayName string     `json:"orgDisplayName"`
	OrgList        []OrgEntry `json:"orgList"`
	KnownRegion    string     `json:"knownRegion"`
}

type OrgEntry struct {
	Name           string `json:"name"`
	DisplayName    string `json:"displayName"`
	SubscriptionId string `json:"subscriptionId"`
}

// Orchestrator types ====================================================
type AppPushRequest struct {
	SandboxId          string
	AppName            string
	DesiredInstances   string
	Overrides          []NVPair
	Filedata           *bytes.Buffer
	EndpointVisibility string
	TunnelAccessKey    string
}

//OrchestratorFTLResponse is response received from FTL Server
type OrchestratorFTLResponse struct {
	Status []string `json:"status"`
}

type OrchestratorResponses struct {
	StatusResponses []OrchestratorResponse `json:"appResponses"`
}

// OrchestratorResponse is the basic response from Orchestrator
type OrchestratorResponse struct {
	AppId     string `json:"appId"`
	Code      string `json:"code"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Details   string `json:"details"`
	Read      bool   `json:"read"`
	LastError string `json:"lastError"`
}

// ReplaceAppInfo struct
type ReplaceAppInfo struct {
	SourceAppID       string      `json:"sourceAppID"`
	TargetAppID       string      `json:"targetAppID"`
	PropertyOverrides []*Property `json:"propertyOverrides"`
}

// END Orchestrator types ====================================================

// DOMAIN SERVER types ====================================================

type DomainServerOrganizationBean struct {
	Id               string `json:"id"`
	OrganizationName string `json:"organizationName"`
	Description      string `json:"description"`
	CreatedTime      int64  `json:"createdTime"`
	LastUpdatedTime  int64  `json:"lastUpdatedTime"`
	LastModifiedBy   string `json:"lastModifiedBy"`
}

// Default Properties
type PropertyDefault struct {
	Name     string `json:"name"`
	DataType string `json:"datatype"`
	Default  string `json:"default"`
}

//  Property  defines a key-value pair to override default values from type "PropertyDefault"
type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Application configuration for an app
type AppConfig struct {
	PropertyPrefix    string             `json:"propertyPrefix"`
	Properties        []*PropertyDefault `json:"properties"`
	PropertyOverrides []*Property        `json:"propertyOverrides"`
}

// DomainServerApplicationBean coming from Domain Server in case of POST or GET of an App
type DomainServerApplicationBean struct {
	Id                           string               `json:"id"`
	ApplicationName              string               `json:"applicationName"`
	Description                  string               `json:"description"`
	Owner                        string               `json:"owner"`
	Version                      string               `json:"version"`
	UserName                     string               `json:"userName"`
	DesiredInstanceCount         uint                 `json:"desiredInstanceCount"`
	CreatedTime                  int64                `json:"createdTime"`
	LastUpdatedTime              int64                `json:"lastUpdatedTime"`
	LastModifiedBy               string               `json:"lastModifiedBy"`
	IsSampleApp                  bool                 `json:"isSampleApp"`
	FullyQualifiedDockerImageUrl string               `json:"fullyQualifiedDockerImageUrl"`
	AppType                      string               `json:"appType,omitempty"`
	EndpointIds                  []string             `json:"endpointIds,omitempty"`
	VolumesFrom                  []string             `json:"volumesFrom"`
	Resources                    *ResourceConstraints `json:"resources"`
	ConfigDetails                *AppConfig           `json:"configDetails"`
	DeploymentStage              string               `json:"deploymentStage"`
	SubscriptionId               string               `json:"subscriptionId"`
	EndpointVisibility           string               `json:"endpointVisibility"`
	IsPackaged                   bool                 `json:"isPackaged"`
	CreatedBy                    string               `json:"createdBy"`
	OwnerName                    string               `json:"ownerName"`
	TibTunnelAccessKey           string               `json:"tibTunnelAccessKey"`
}

// DomainServerApplicationsResponse coming from Domain Server for GET apps
type DomainServerApplicationsResponse struct {
	ApplicationBeans []DomainServerApplicationBean `json:"applications"`
}

// DomainServerAppEndpointUrlResponse coming from Domain Server for app endpoint url based on the given endpointId in request
type DomainServerAppEndpointUrlResponse struct {
	EndpointUrl string `json:"endpointUrl"`
}

// DomainServerAppEndpointBean is the response from the Domain Server for GET app endpoint request
type DomainServerAppEndpointBean struct {
	// NOTE: Unused JSON elements are intentionally omitted.
	Type string `json:"type"` // the endoint type: public or private
}

// DomainServerGetSandboxesResponse is response from Domain Server for GET sandboxes request
type DomainServerGetSandboxesResponse struct {
	Sandboxes        []DomainServerSandboxBean `json:"sandboxes"`
	OperationWarning string                    `json:"operationWarning"`
}

// DomainServerSandboxBean is response from Domain Server for GET sandbox details request
type DomainServerSandboxBean struct {
	Id              string   `json:"id"`
	SandboxName     string   `json:"sandboxName"`
	DisplayName     string   `json:"displayName"`
	Description     string   `json:"description"`
	OrganizationId  string   `json:"organizationId"`
	CreatedBy       string   `json:"createdBy"`
	UsedBy          []string `json:"usedBy"`
	ApplicationIds  []string `json:"applicationIds"`
	CreatedTime     int64    `json:"createdTime"`
	LastUpdatedTime int64    `json:"lastUpdatedTime"`
	LastModifiedBy  string   `json:"lastModifiedBy"`
	SandboxType     string   `json:"sandboxType"`
	Visibility      string   `json:"visibility"`
	EndpointType    string   `json:"endpointType"`
}

type DomainServerUserBean struct {
	FirstName        string   `json:"firstName"`
	LastName         string   `json:"lastName"`
	UserName         string   `json:"userName"`
	Email            string   `json:"email"`
	CompanyName      string   `json:"companyName"`
	Phone            string   `json:"phone"`
	UserId           string   `json:"userId"`
	DefaultSandboxId string   `json:"defaultSandboxId"`
	SandboxIds       []string `json:"sandboxIds"`
	OrganizationId   string   `json:"organizationId"`
	RepoUrl          string   `json:"repoUrl"`
	Status           bool     `json:"status"`
	LastUpdatedTime  int64    `json:"lastUpdatedTime"`
	EulaAcceptedTime int64    `json:"eulaAcceptedTime"`
	Disabled         bool     `json:"disabled"`
	SubscriptionKey  string   `json:"subscriptionKey"`
}

type DomainServerAppAudits struct {
	TotalNum          int                    `json:"totalNum"`
	LastEvaluatedTime string                 `json:"queryLocator"`
	Audits            []DomainServerAppAudit `json:"audits"`
}

type DomainServerAppAudit struct {
	AppId                string `json:"appId"`
	CreatedTime          int64  `json:"createdTime"`
	UserId               string `json:"userId"`
	UserName             string `json:"userName"`
	Action               string `json:"action"`
	ActionSummary        string `json:"actionSummary"`
	StatusCode           string `json:"statusCode"`
	Duration             int    `json:"duration"`
	PushPrepDuration     int    `json:"pushPrepDuration"`
	PushValidateDuration int    `json:"pushValidateDuration"`
	PushBuildDuration    int    `json:"pushBuildDuration"`
	PushScaleDuration    int    `json:"pushScaleDuration"`
	AppType              string `json:"appType"`
	Gsbc                 string `json:"gsbc"`
	Client               string `json:"client"`
}

// END DOMAIN SERVER types ==================================================

// APP MANAGER types =========================================================

// AppManagerGetAppInfoResponse is response from app manager for get app information request
type AppManagerGetAppInfoResponse struct {
	Id        string `json:"id"`        //appId
	Instances uint   `json:"instances"` //healthyInstanceCount
}

// END APP MANAGER types ======================================================

// ErrorResponse is error response from domain server, build server, app manager , orchestrator
type ErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	// ErrorDetail string `json:"errorDetail, omitempty"`
	ErrorDetail string `json:"errorDetail"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

// Manifest defines the struct of app manifest file
type Manifest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

//Responses from OAuth2 port of TIBCO Accounts
type OAResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	Error        string `json:"error"`
	ErrorDesc    string `json:"error_description"`
	Revoked      string `json:"revoked"`
}

//for requests with username-password to TA
type AuthRequest struct {
	Username,
	Pwd,
	ClientId string
}

//for followup requests: renewal & logout
type FollowupRequest struct {
	Token,
	ClientId,
	ClientSecret string
}

// ResourceConstraints contains resource constraints for an application instance
type ResourceConstraints struct {
	PhysicalMemory uint `json:"physicalMemory"`
	SwapMemory     uint `json:"swapMemory"`
	CpuQuota       uint `json:"cpuQuota"`
}

// NVPair is a type for generic name/value pairs as used to retrieve/set configuration
type NVPair struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

type OrgFTLStatus struct {
	IsFTLEnabled bool `json:"isFTLEnabled"`
}

type TunnelAction struct {
	AccessKey string
	Action    string
}

type OrganizationInfo struct {
	Gsbc      string `json:"gsbc"`
	AppDomain string `json:"appDomain"`
}

type BuildServerResponse struct {
	Code    string `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type RestCallRequest struct {
	Method       string                  // REST method
	Url          *url.URL                // URL of the Rest API to be invoked
	Headers      map[string]string       // HTTP Headers to be passed while invoking the API
	Body         io.Reader               // Request body for the API
	LogRequest   bool                    // if true; then the rest handler will log the request otherwise will skip it
	UserId       string                  // user Id used for logging
	RetryAttempt *utilities.RetryAttempt // retry times if connection failed
}
