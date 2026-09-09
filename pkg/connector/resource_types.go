package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

// userResourceType is for all user objects from App Store Connect.
var userResourceType = &v2.ResourceType{
	Id:          "user",
	DisplayName: "User",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_USER},
}

// appResourceType is for all app objects from App Store Connect.
var appResourceType = &v2.ResourceType{
	Id:          "app",
	DisplayName: "App",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_APP},
}

// roleResourceType is a synthetic resource that represents the App Store Connect
// account itself. Roles are entitlements on this resource.
var roleResourceType = &v2.ResourceType{
	Id:          "role",
	DisplayName: "Role",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_GROUP},
}

// betaGroupResourceType is for TestFlight beta groups. Membership of a beta
// group is what grants a tester access to an app's builds.
var betaGroupResourceType = &v2.ResourceType{
	Id:          "beta_group",
	DisplayName: "Beta Group",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_GROUP},
}

// betaTesterResourceType is for TestFlight beta testers. These are distinct
// from App Store Connect users: a tester is an email address invited to test
// builds, and may not have an App Store Connect account at all.
var betaTesterResourceType = &v2.ResourceType{
	Id:          "beta_tester",
	DisplayName: "Beta Tester",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_USER},
}
