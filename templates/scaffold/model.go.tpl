package model

// TODO: implement another field
type {{.ModelName}} struct {
	BaseAssociatedEntity
	Name string `json:"name" firestore:"name"`
}

type {{.ModelName}}ListOption struct {
	UserID         string
	OrganizationID string
}

