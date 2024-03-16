# service について

ある model に関して操作を行うための関数を持つ

基本的には 該当 model に関する操作のみを行い、他 Service を呼び出すような処理は usecase に書く

## 基盤として実装されているもの

ベースとなる Interface は以下のような設計になっている

```go
type AbstractGenericService[T model.AbstractEntity, U model.AbstractListOption] interface {
	Get(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object T) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	Set(ctx context.Context, id string, object T) error
	List(ctx context.Context, opts U) ([]T, error)
}

```

この Interface を満たしていると、このまま usecase.baseGenericUseCase で使えるようになる

`AbstractEntity`, `AbstractListOption` については[model の詳細ページ参照](./model.md)

### 例

```go
type NoteService interface {
	AbstractGenericService[*model.Note, *model.NoteListOption]
}

type noteService struct {
	repo repository.NoteRepository
	baseGenericService[*model.Note, *model.NoteListOption]
}

func NewNoteService(repo repository.NoteRepository) NoteService {
	base := newGenericService(repo)
	return &noteService{repo: repo, baseGenericService: base}
}

var _ NoteService = &noteService{}

```

`baseGenericService` 側にて必要な関数が定義しされているので、上記だけで`AbstractGenericService`を満たす

## 拡張

```go

type TeamService interface {
	AbstractGenericService[*model.Team, *model.TeamListOption]

	AddMember(ctx context.Context, teamID string, user *model.User) error
	DeleteMember(ctx context.Context, teamID string, user *model.User) error
}

type teamService struct {
	repo repository.TeamRepository
	baseGenericService[*model.Team, *model.TeamListOption]
}

func NewTeamService(repo repository.TeamRepository) TeamService {
	base := newGenericService(repo)
	return &teamService{repo: repo, baseGenericService: base}
}

func (s *teamService) AddMember(ctx context.Context, teamID string, user *model.User) error {
	return s.repo.AddMember(ctx, teamID, user)
}

func (s *teamService) DeleteMember(ctx context.Context, teamID string, user *model.User) error {
	return s.repo.DeleteMember(ctx, teamID, user)
}

```

このように interface に変更がなければ新たな関数(`AddMember`)を追加することも可能

## Update について

- firestore 以外の repository に対応するため、`firestore.Increment`など特殊なパターンに対応できない
- 対応するべく、代わりに`model.IncrementOperation` `model.ArrayUnionOperation` `model.RemoveUnionOperation`を定義しているのでそれを使う。
  - 詳しくは実装を参照。
