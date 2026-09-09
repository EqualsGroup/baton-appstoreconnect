package client

// JSON:API envelope types for App Store Connect API responses.

// Response is the top-level JSON:API response envelope.
type Response[T any] struct {
	Data  []T             `json:"data"`
	Links PaginationLinks `json:"links,omitempty"`
}

// SingleResponse is a JSON:API response with a single data object.
type SingleResponse[T any] struct {
	Data  T               `json:"data"`
	Links PaginationLinks `json:"links,omitempty"`
}

// PaginationLinks contains pagination URLs from the App Store Connect API.
type PaginationLinks struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
}

// User represents an App Store Connect user resource.
type User struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Attributes UserAttributes `json:"attributes"`
}

// UserAttributes contains the attributes of a user.
type UserAttributes struct {
	Username            string   `json:"username"`
	FirstName           string   `json:"firstName"`
	LastName            string   `json:"lastName"`
	Email               string   `json:"email"`
	Roles               []string `json:"roles"`
	AllAppsVisible      bool     `json:"allAppsVisible"`
	ProvisioningAllowed bool     `json:"provisioningAllowed"`
}

// App represents an App Store Connect app resource.
type App struct {
	Type       string        `json:"type"`
	ID         string        `json:"id"`
	Attributes AppAttributes `json:"attributes"`
}

// AppAttributes contains the attributes of an app.
type AppAttributes struct {
	Name     string `json:"name"`
	BundleID string `json:"bundleId"`
	SKU      string `json:"sku"`
}

// UserInvitationCreateRequest is the request body for inviting a new user.
type UserInvitationCreateRequest struct {
	Data UserInvitationCreateData `json:"data"`
}

// UserInvitationCreateData is the data object for creating a user invitation.
type UserInvitationCreateData struct {
	Type          string                         `json:"type"`
	Attributes    UserInvitationCreateAttributes `json:"attributes"`
	Relationships *UserInvitationRelationships   `json:"relationships,omitempty"`
}

// UserInvitationCreateAttributes contains the attributes for a user invitation.
type UserInvitationCreateAttributes struct {
	Email          string   `json:"email"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	Roles          []string `json:"roles"`
	AllAppsVisible bool     `json:"allAppsVisible"`
}

// UserInvitationRelationships contains optional relationships for user invitations.
type UserInvitationRelationships struct {
	VisibleApps *VisibleAppsRelationship `json:"visibleApps,omitempty"`
}

// VisibleAppsRelationship contains the list of visible apps for a user invitation.
type VisibleAppsRelationship struct {
	Data []RelationshipData `json:"data"`
}

// RelationshipData is a generic JSON:API relationship data item.
type RelationshipData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// UserVisibleAppsResponse is the response for listing a user's visible apps.
type UserVisibleAppsResponse struct {
	Data  []App           `json:"data"`
	Links PaginationLinks `json:"links,omitempty"`
}

// AppStoreConnectRole defines known App Store Connect roles.
var AppStoreConnectRoles = []string{
	"ADMIN",
	"FINANCE",
	"ACCOUNT_HOLDER",
	"SALES",
	"DEVELOPER",
	"APP_MANAGER",
	"CUSTOMER_SUPPORT",
	"MARKETING",
	"CREATE_APPS",
	"CLOUD_MANAGED_DEVELOPER_ID",
	"CLOUD_MANAGED_APP_DISTRIBUTION",
	"ACCESS_TO_REPORTS",
	"GENERATE_INDIVIDUAL_KEYS",
	"IMAGE_MANAGER",
}

// BetaGroup represents a TestFlight beta group resource.
type BetaGroup struct {
	Type          string                  `json:"type"`
	ID            string                  `json:"id"`
	Attributes    BetaGroupAttributes     `json:"attributes"`
	Relationships *BetaGroupRelationships `json:"relationships,omitempty"`
}

// BetaGroupAttributes contains the attributes of a TestFlight beta group.
type BetaGroupAttributes struct {
	Name                   string `json:"name"`
	CreatedDate            string `json:"createdDate"`
	IsInternalGroup        bool   `json:"isInternalGroup"`
	HasAccessToAllBuilds   bool   `json:"hasAccessToAllBuilds"`
	PublicLinkEnabled      bool   `json:"publicLinkEnabled"`
	PublicLink             string `json:"publicLink"`
	PublicLinkLimitEnabled bool   `json:"publicLinkLimitEnabled"`
	PublicLinkLimit        int    `json:"publicLinkLimit"`
	FeedbackEnabled        bool   `json:"feedbackEnabled"`
}

// BetaGroupRelationships contains the relationships of a beta group. Only the
// owning app is requested by this connector.
type BetaGroupRelationships struct {
	App *SingleRelationship `json:"app,omitempty"`
}

// AppID returns the ID of the app that owns this beta group, or an empty string
// when the relationship was not included in the response.
func (b BetaGroup) AppID() string {
	if b.Relationships == nil || b.Relationships.App == nil || b.Relationships.App.Data == nil {
		return ""
	}

	return b.Relationships.App.Data.ID
}

// BetaTester represents a TestFlight beta tester resource.
type BetaTester struct {
	Type       string               `json:"type"`
	ID         string               `json:"id"`
	Attributes BetaTesterAttributes `json:"attributes"`
}

// BetaTesterAttributes contains the attributes of a TestFlight beta tester.
type BetaTesterAttributes struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	InviteType string `json:"inviteType"`
	State      string `json:"state"`
}

// BetaTesterCreateRequest is the request body for creating a beta tester.
type BetaTesterCreateRequest struct {
	Data BetaTesterCreateData `json:"data"`
}

// BetaTesterCreateData is the data object for creating a beta tester.
type BetaTesterCreateData struct {
	Type          string                     `json:"type"`
	Attributes    BetaTesterCreateAttributes `json:"attributes"`
	Relationships *BetaTesterRelationships   `json:"relationships,omitempty"`
}

// BetaTesterCreateAttributes contains the attributes for creating a beta tester.
type BetaTesterCreateAttributes struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

// BetaTesterRelationships contains optional relationships for a beta tester.
type BetaTesterRelationships struct {
	BetaGroups *RelationshipRequest `json:"betaGroups,omitempty"`
}

// RelationshipRequest is a JSON:API to-many relationship body, used both for
// nested relationships and for relationship endpoint requests.
type RelationshipRequest struct {
	Data []RelationshipData `json:"data"`
}

// SingleRelationship is a JSON:API to-one relationship.
type SingleRelationship struct {
	Data *RelationshipData `json:"data,omitempty"`
}
