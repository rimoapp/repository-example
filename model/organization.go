package model

// TODO: implement another field
type Organization struct {
	BaseAssociatedEntity
	Name string `json:"name" firestore:"name"`
}

type OrganizationListOption struct {
	UserID         string
	OrganizationID string
}
