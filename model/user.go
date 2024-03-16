package model

// TODO: implement another field
type User struct {
	BaseAssociatedEntity
	Name  string `json:"name" firestore:"name"`
	Teams []Team `json:"teams,omitempty" firestore:"-" gorm:"many2many:team_members;"`
}

type UserListOption struct {
	BaseListOption
}
