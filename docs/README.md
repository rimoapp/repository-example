# この github repository について

- `repository-example` という github repository について説明します
- この `repository-example` は、Go 言語での API サーバーの実装で、repository 層を用いた設計を行っています

## 設計方針

- model 層 / repository 層 / service 層 / handler 層 と大きく 4 つの層に分かれる
- 現状 firestore と gorm を用いた RDB に対応している

### model

- 各 model ごとに、model 名や Field 等を定義する
  - ex. Note, User ..
- また、response の時の json 名、firestore に保存するカラム名、gorm に関する設定もここで定義する
- [詳細ページ](./model.md)

### repository

- DB へのアクセスに関するロジックを持つ
- 各 model ごとに Repository が作られれるが、基本的に必要な関数は共通化されている
  - ex. NoteRepository, UserRepository ..
  - Create / Get / Update / Set / Delete / List 関数を持つ
  - もちろん repository によっては使わない関数も存在する
- 現状 firestore 用と gorm 経由の RDB 用の 2 つの repository がある
- [詳細ページ](./repository.md)

### Service

- 各 model ごとに service が作られる
  - ex. NoteService, UserService ..
- repository 層を受け取って生成される
- 基本的に Create / Get / Update / Set / Delete / List 関数を持ち、それぞれ内部で repository の関数を呼び出す
  - もちろん service によっては使わない関数も存在する
- [詳細ページ](./service.md)

### Handler

- 各 model ごとに handler が作られる
  - ex. NoteHandler, UserHandler ..
- service 層を受け取って生成される
- [handler の詳細ページ](./handler.md)

## メリット

- repository 層に DB アクセスのロジックを集約させることで、他の箇所との結合が緩和される
  - また、DB アクセスのロジックを変更する場合、repository 層のみを変更すれば良い
  - DB を切り替える際も repository 層を切り替えれば良い
- 各 model ごとに共通可能なところを共通にすることで実装の効率が図れる。
- また、scaffold(後述)的に簡単にベースが実装可能
  - ついでに test も書きやすくなったはず

## scaffold

- 今回整理した関係で共通化できるところが増えたので scaffold cmd を作ったのでこれで土台ができる
  - 新たにモデルを作るときに呼び出す想定
- ex: `go run ./cmd/scaffold/main.go Note`
  - `Note`の部分は CamelCase で新モデル名を指定
- 生成されるファイルは以下
  - model/note.go
  - service/note.go
  - service/note_test.go
  - handler/note.go
  - repository/note.go
- メモ
  - model の各 Field は実装されないので cmd 実行後実装してください
  - handler.go に SetRouter が実装されますが、それは router.go から呼び出せます
  - 例
    - handler/note.go
      ```go
        func (h *NoteHandler) SetRouter(group *gin.RouterGroup) {
          group.GET("/notes", h.List)
          group.GET("/notes/:ID", h.Get)
          group.DELETE("/notes/:ID", h.Delete)
          group.PATCH("/notes/:ID", h.Update)
          group.POST("/notes", h.Create)
        }
      ```
  - router.go
    ```go
      ...
      noteHandler.SetRouter(public)
      ...
    ```
