package consts

const (
	CLI_MODULE_NAME       string = "tas-cli"
	CLI_MODULE_USAGE      string = "TIBCO AuditSafe Cloud - Command Line Interface"
	CLI_VERSION           string = "1.0.0"
	CLI_COPYRIGHT_MESSAGE string = "2015-2019 TIBCO Software Inc."

	//environment property for indicating extra debug info; sensitive data will be displayed. Don't use
	TASCLI_DBG string = "TASCLI_DBG"

	OBFUSCATE_COOKIE_VALUE = true

	//environment property governing persistence of OAuth tokens
	DONT_PERSIST string = "TIBCLI_DONT_PERSIST"
)
