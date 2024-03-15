package model

// TODO: implement another field
type Organization struct {
	BaseAssociatedEntity
	Name  string  `json:"name" firestore:"name"`
	Teams []*Team `json:"teams,omitempty" firestore:"-"`
}

type OrganizationListOption struct {
	BaseListOption
}

type GetOrganizationOption struct {
	IncludeTeams bool `form:"include_teams" json:"include_teams" binding:"omitempty"`
}
