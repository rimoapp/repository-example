package model

type Team struct {
	BaseAssociatedEntity
	Name           string  `json:"name" firestore:"name"`
	OrganizationID string  `json:"organization_id" firestore:"organization_id"`
	Members        []*User `json:"members,omitempty" firestore:"-" gorm:"many2many:team_members;"`
}

type TeamListOption struct {
	BaseListOption
	OrganizationID string
}
