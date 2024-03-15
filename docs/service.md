# service について

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

この Interface を満たしていると、このまま BaseGenericHandler で使えるようになる

`AbstractEntity`, `AbstractListOption` については[model の詳細ページ参照](./model.md)

### 例

```go
type NoteService struct {
	Repo repository.NoteRepository
	baseGenericService[*model.Note, *model.NoteListOption]
}

func NewNoteService(repo repository.NoteRepository) *NoteService {
	base := newGenericService(repo)
	return &NoteService{Repo: repo, baseGenericService: base}
}

```

`baseGenericService` 側にて必要な関数が定義しされているので、上記だけで`AbstractGenericService`を満たす

## 拡張

```go

type OrganizationService struct {
	Repo repository.OrganizationRepository
	baseGenericService[*model.Organization, *model.OrganizationListOption]
	TeamService *TeamService
}

func NewOrganizationService(repo repository.OrganizationRepository, teamService *TeamService) *OrganizationService {
	base := newGenericService(repo)
	return &OrganizationService{Repo: repo, baseGenericService: base, TeamService: teamService}
}

func (s *OrganizationService) Create(ctx context.Context, object *model.Organization) (string, error) {
	now := time.Now()
	object.CreatedAt = now
	object.UpdatedAt = now
	return s.Repo.Create(ctx, object)
}

func (s *OrganizationService) GetWithOption(ctx context.Context, id string, opts *model.GetOrganizationOption) (*model.Organization, error) {
	org, err := s.Repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if opts != nil && opts.IncludeTeams {
		teams, err := s.TeamService.List(ctx, &model.TeamListOption{OrganizationID: id})
		if err != nil {
			return nil, err
		}
		org.Teams = teams
	}
	return org, nil
}

```

このように interface に変更がなければ既存の関数(`Create`)も再定義しても問題ない。新たな関数(`GetWithOption`)を追加することも可能

また、`TeamService` のように他の service を埋め込むことも可能

## Update について

- firestore 以外の repository に対応するため、`firestore.Increment`など特殊なパターンに対応できない
- 対応するべく、代わりに`model.IncrementOperation` `model.ArrayUnionOperation` `model.RemoveUnionOperation`を定義しているのでそれを使う。
  - 詳しくは実装を参照。
