# repository について

DB にアクセスするロジックを集約する層で、基本的に DB に関する処理は全てここで行う。

現状 firestore と gorm を用いた RDB に対応している。

また model ごとに repository を持ち、interface としては共通にしている

## 例

NoteRepository は Get / Create / Update / Set / Delete / List の関数を持つ

また、 List 以外は firestore / gorm ともに関数自体は共通にしやすいが、List に関する where 句の生成は firestore / gorm で異なり、かつ ListOption ごとにも異なるため、List は firestore / gorm で別々に実装する必要がある

```go

type NoteRepository interface {
	Get(ctx context.Context, id string) (*model.Note, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object *model.Note) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	List(ctx context.Context, opts *model.NoteListOption) ([]*model.Note, error)
	Set(ctx context.Context, id string, object *model.Note) error
}

const notesCollectionPath = "notes"

func NewNoteRepository(opts NewRepositoryOption) NoteRepository {
	if opts.dbClient != nil {
		return &gormNoteRepository{
			gormGenericRepository: *newGormGenericRepository[*model.Note](opts.dbClient),
		}
	}
	return &firestoreNoteRepository{
		firestoreGenericRepository: *newFirestoreGenericRepository[*model.Note](opts.firestoreClient, notesCollectionPath),
	}
}

type firestoreNoteRepository struct {
	firestoreGenericRepository[*model.Note]
}

func (r *firestoreNoteRepository) List(ctx context.Context, opts *model.NoteListOption) ([]*model.Note, error) {
	collectionPath := r.collectionPath
	if collectionPath == "" {
		return nil, errors.New("collection path is empty")
	}
	query := r.client.Collection(collectionPath).Query
	if opts.UserID != "" {
		query = query.Where("user_id", "==", opts.UserID)
	}
	return r.firestoreGenericRepository.list(ctx, query)
}

type gormNoteRepository struct {
	gormGenericRepository[*model.Note]
}

func (r *gormNoteRepository) List(ctx context.Context, opts *model.NoteListOption) ([]*model.Note, error) {
	query := r.client
	if opts.UserID != "" {
		query = query.Where("user_id = ?", opts.UserID)
	}
	return r.gormGenericRepository.list(query)
}

```

## 拡張

上記 Interface は満たさないと GenericService や GenericHandler が使えなくなるので注意が必要

ただ関数が増える分には問題はないので、必要な関数は追加することで対応可能

```go
type TeamRepository interface {
	Get(ctx context.Context, id string) (*model.Team, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object *model.Team) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	List(ctx context.Context, opts *model.TeamListOption) ([]*model.Team, error)
	Set(ctx context.Context, id string, object *model.Team) error
	AddMember(ctx context.Context, teamID string, user *model.User) error
	DeleteMember(ctx context.Context, teamID string, user *model.User) error
}

```
