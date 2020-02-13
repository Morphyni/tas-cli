package types

//OrgInfo to store account name for multi subscription login request
type OrgInfo struct {
	AccountName string
	Region      string
}
type OrgEntry struct {
	Name           string `json:"name"`
	DisplayName    string `json:"displayName"`
	SubscriptionId string `json:"subscriptionId"`
}
