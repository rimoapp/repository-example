# model について

repository 層にて扱う関係で model の制約について記載する。

## Model

### AbstractEntity

```go
type AbstractEntity interface {
	GetID() string
	SetID(id string)
}
```

repository 層にて扱うモデルは上記 `AbstractEntity` を満たす

これを満たした基底モデルとして以下の `BaseEntity` が定義される

```go
type BaseEntity struct {
	ID        string    `json:"id" firestore:"-"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`

	// for gorm
	UID       int            `json:"-" firestore:"-" gorm:"primarykey"`
	DeletedAt gorm.DeletedAt `json:"-" firestore:"-" gorm:"index"`
}
```

ので、repository 層で扱うモデルは以下のように`BaseEntity`を埋め込むことで実装してほしい

```go
type Note struct {
    BaseEntity
    ...(other fields)
}
```

### AbstractAssociatedEntity

```go
type AbstractAssociatedEntity interface {
	AbstractEntity
	IsAuthorized(userID string) bool
	SetCreatorID(userID string)
}
```

handler 層にて扱うモデルは上記`AbstractAssociatedEntity` を満たす

`AbstractEntity` を拡張させたもので、handler 層にてログインユーザを用いた権限チェックを可能にしている

これを満たした基底モデルとして以下の `BaseAssociatedEntity` が定義している

```go
type BaseAssociatedEntity struct {
	BaseEntity
	UserID string `json:"user_id" firestore:"user_id"`
}
```

ので、handler 層で扱うモデルは以下のように`BaseAssociatedEntity`を埋め込むことで実装してほしい

```go
type Note struct {
    BaseAssociatedEntity
    ...(other fields)
}
```

## ListOption

各 Model ごとに ListOption を定義している。ListOption は以下の interface, `AbstractListOption`を満たす

### AbstractListOption / BaseListOption

```go
type AbstractListOption interface {
	SetUserID(userID string)
}
```

`AbstractListOption`は repository 層にて扱う List 関数にて、where 句を生成する際に用いられる interface.

handler 側でログインユーザ ID をセットさせるため、SetUserID を定義している

AbstractListOption を満たす BaseListOption も併せて定義しているので、各 model の ListOption は以下のように`BaseListOption`を埋め込んで実装する

```go
type BaseListOption struct {
	UserID string
}
func (b *BaseListOption) SetUserID(userID string) {
	b.UserID = userID
}
var _ AbstractListOption = (*BaseListOption)(nil)

```

```go
type NoteListOption struct {
	BaseListOption
    ...(other fields)
}
```

## AbstractUpdateOperation

repository 層を分離するにあたって、`firestore.ArrayUnion`等をそのまま用いることができない

そのため `AbstractUpdateOperation`を定義しているのでそちらを用いて更新作業を行う

詳しくは実装を参照。
