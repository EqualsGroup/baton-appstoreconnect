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
