# model について

repository 層にて扱う関係で model の制約について記載します。

## Model

### AbstractDataEntity

```go
type AbstractDataEntity interface {
	GetID() string
	SetID(id string)
}
```

repository 層にて扱うモデルは上記 `AbstractDataEntity` を満たしてください

これを満たした基底モデルとして以下の `BaseDataEntity` が定義されています

```go
type BaseDataEntity struct {
	ID        string    `json:"id" firestore:"-"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`

	// for gorm
	UID       int            `json:"-" firestore:"-" gorm:"primarykey"`
	DeletedAt gorm.DeletedAt `json:"-" firestore:"-" gorm:"index"`
}
```

ので、repository 層で扱うモデルは以下のように`BaseDataEntity`を埋め込むことで実装してください

```go
type Note struct {
    BaseDataEntity
    ...(other fields)
}
```

### AbstractAssociatedEntity

```go
type AbstractAssociatedEntity interface {
	AbstractDataEntity
	IsAuthorized(userID string) bool
	SetCreatorID(userID string)
}
```

handler 層にて扱うモデルは上記`AbstractAssociatedEntity` を満たしてください

`AbstractDataEntity` を拡張させたもので、handler 層にてログインユーザを用いた権限チェックを可能にしています

これを満たした基底モデルとして以下の `BaseAssociatedEntity` が定義されています

```go
type BaseAssociatedEntity struct {
	BaseDataEntity
	UserID string `json:"user_id" firestore:"user_id"`
}
```

ので、handler 層で扱うモデルは以下のように`BaseAssociatedEntity`を埋め込むことで実装してください

```go
type Note struct {
    BaseAssociatedEntity
    ...(other fields)
}
```

## ListOption

各 Model ごとに ListOption を持ちます。ListOption は以下の interface, `AbstractListOption`を満たしてください

### AbstractListOption / BaseListOption

```go
type AbstractListOption interface {
	SetUserID(userID string)
}
type BaseListOption struct {
	UserID string
}
func (b *BaseListOption) SetUserID(userID string) {
	b.UserID = userID
}
var _ AbstractListOption = (*BaseListOption)(nil)
```

repository 層にて扱う List 関数にて、where 句を生成する際に用いられる interface です

handler 側でログイン ID をセットさせるため、SetUserID を定義しています

AbstractListOption を満たす BaseListOption も併せて定義しているので、各 model の ListOption は以下のように`BaseListOption`を埋め込むことで実装してください

```go
type NoteListOption struct {
	BaseListOption
    ...(other fields)
}
```

## AbstractUpdateOperation

repository 層を分離するにあたって、`firestore.ArrayUnion`等をそのまま用いることができません

そのため `AbstractUpdateOperation`を定義しています

<!-- ちゃんと書く -->
