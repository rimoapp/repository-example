package model

// TODO: implement another field
type User struct {
	BaseAssociatedEntity
	Name string `json:"name" firestore:"name"`
}

type UserListOption struct {
	BaseListOption
}
