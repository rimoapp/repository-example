package model

// TODO: implement another field
type Team struct {
	BaseAssociatedEntity
	Name           string `json:"name" firestore:"name"`
	OrganizationID string `json:"organization_id" firestore:"organization_id"`
}

type TeamListOption struct {
	UserID string
}
