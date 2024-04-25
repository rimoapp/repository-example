# usecase

service を受け取り、その service を呼び出すような処理を行う

もし他の service を呼び出すような処理がある場合、その記述をここに書く

## 基盤として実装されているもの

ベースとなる Interface は以下のような設計になっている

```go
type AbstractGenericUseCase[T model.AbstractEntity, U model.AbstractListOption] interface {
	Get(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object T) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	Set(ctx context.Context, id string, object T) error
	List(ctx context.Context, opts U) ([]T, error)
}
```

この Interface を満たしていると、このまま handler.baseGenericHandler で使えるようになる

`AbstractEntity`, `AbstractListOption` については[model の詳細ページ参照](./model.md)

### 例

```go
type NoteUseCase interface {
	AbstractGenericUseCase[*model.Note, *model.NoteListOption]
}

type noteUseCase struct {
	baseGenericUseCase[*model.Note, *model.NoteListOption]
    svc service.NoteService
}

func NewNoteUseCase(svc service.NoteService) NoteUseCase {
	base := newGenericUseCase(svc)
	return &noteUseCase{baseGenericUseCase: base, svc:svc}
}

var _ NoteUseCase = &noteUseCase{}

```

`baseGenericUseCase` 側にて必要な関数が定義しされているので、上記だけで`AbstractGenericUseCase`を満たす

## 拡張

```go

type OrganizationUseCase interface {
	AbstractGenericUseCase[*model.Organization, *model.OrganizationListOption]
	GetWithOption(ctx context.Context, id string, opts *model.GetOrganizationOption) (*model.Organization, error)
}

type organizationUseCase struct {
	baseGenericUseCase[*model.Organization, *model.OrganizationListOption]
	svc         service.OrganizationService
	teamService service.TeamService
}

func NewOrganizationUseCase(svc service.OrganizationService, teamService service.TeamService) OrganizationUseCase {
	base := newGenericUseCase(svc)
	return &organizationUseCase{baseGenericUseCase: base, svc: svc, teamService: teamService}
}

var _ OrganizationUseCase = &organizationUseCase{}

func (uc *organizationUseCase) GetWithOption(ctx context.Context, id string, opts *model.GetOrganizationOption) (*model.Organization, error) {
	org, err := uc.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if opts != nil && opts.IncludeTeams {
		teams, err := uc.teamService.List(ctx, &model.TeamListOption{OrganizationID: id})
		if err != nil {
			return nil, err
		}
		org.Teams = teams
	}
	return org, nil
}

```

上記例は、`AbstractGenericUseCase` を満たすと同時に、`GetWithOption` という関数を追加している

また、`GetWithOption` にて他の service を呼び出している
