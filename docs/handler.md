# handler について

handler は、HTTP リクエストを受け取り、service 層を呼び出し、レスポンスを返す役割を持つ。

基底クラスとして`baseGenericHandler[T model.AbstractAssociatedEntity, U model.AbstractListOption]` を定義しており、これを満たすように実装されていると、以下のような関数がそのまま使える。

```go
func (h *baseGenericHandler[T, U]) Get(c *gin.Context)
func (h *baseGenericHandler[T, U]) Delete(c *gin.Context)
func (h *baseGenericHandler[T, U]) Create(c *gin.Context)
func (h *baseGenericHandler[T, U]) Update(c *gin.Context)
func (h *baseGenericHandler[T, U]) List(c *gin.Context)
```

## 例

```go
type NoteHandler struct {
	baseGenericHandler[*model.Note, *model.NoteListOption]
	svc *service.NoteService
}

func NewNoteHandler(opts repository.NewRepositoryOption) *NoteHandler {
	repo := repository.NewNoteRepository(opts)
	svc := service.NewNoteService(repo)
	handler := NewGenericHandler(svc, "noteID")
	return &NoteHandler{baseGenericHandler: *handler, svc: svc}
}

func (h *NoteHandler) SetRouter(group *gin.RouterGroup) {
	group.GET("/notes", h.List)
	group.GET("/notes/:noteID", h.Get)
	group.POST("/notes", h.Create)
	group.DELETE("/notes/:noteID", h.Delete)
	group.PATCH("/notes/:noteID", h.Update)
}

```

これだけで、`Note` に対する CRUD API が実装される。

また、scaffold cmd を用いて実装すると`SetRouter` が自動で実装されるので、以下のように router から呼び出すだけで使える。

```go
    ...
	opts, err := repository.BuildNewRepositoryOptions(ctx)
    ...
	noteHandler := handler.NewNoteHandler(*opts)
	noteHandler.SetRouter(rootGroup)
    ...
```

## 拡張

```go
func (h *OrganizationHandler) Get(c *gin.Context) {
	req := &model.GetOrganizationOption{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}

	id := c.Param(h.idParam)
	entity, err := h.svc.GetWithOption(c, id, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to get").Error()})
		return
	}
	if !entity.IsAuthorized(c.GetString("user_id")) {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	c.JSON(http.StatusOK, entity)
}
```

このように例えば`Get`関数を拡張することで、service 層に独自に追加した関数を handler 側から呼び出すことができる。
