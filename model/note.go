package model

// TODO: implement another field
type Note struct {
	BaseAssociatedEntity
	Name string `json:"name" firestore:"name"`
}

type NoteListOption struct {
	UserID         string
}
