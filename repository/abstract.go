package repository

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"google.golang.org/api/option"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AbstractGenericRepository[T model.AbstractDataEntity] interface {
	Get(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object T) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	Set(ctx context.Context, id string, object T) error
	// List(ctx context.Context, options interface{}) ([]T, error)
}

type NewRepositoryOption struct {
	FirestoreClient *firestore.Client
	DBClient        *gorm.DB
}

func BuildNewRepositoryOptions(ctx context.Context) (*NewRepositoryOption, error) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") != "" {
		fmt.Println("firestore repository")
		conf := &firebase.Config{ProjectID: "your-project-id"} // プロジェクトIDを設定
		app, err := firebase.NewApp(ctx, conf, option.WithoutAuthentication())
		if err != nil {
			return nil, fmt.Errorf("error initializing app: %v", err)
		}
		client, err := app.Firestore(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to initialize firestore client")
		}
		return &NewRepositoryOption{
			FirestoreClient: client,
		}, nil
	}
	fmt.Println("gorm repository")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &NewRepositoryOption{
		DBClient: db,
	}, nil
}
func createNewInstance[T model.AbstractDataEntity]() T {
	var entity T
	// reflect.TypeOf を使用して T の型情報を取得
	t := reflect.TypeOf(entity)
	// T がポインタ型の場合、新しいインスタンスを生成
	if t.Kind() == reflect.Ptr {
		// 新しいインスタンスを生成して返す
		newInstance := reflect.New(t.Elem()).Interface()
		return newInstance.(T)
	}
	// ポインタでない場合はそのままデフォルト値を返す
	return entity
}

func validateKeyValues(keyValues map[string]interface{}) error {
	for _, v := range keyValues {
		switch v.(type) {
		// 許可する型
		case string, int, float64, time.Time, int32, bool, []string, []int, []int32, []bool, []float64, map[string]int32, map[string]string, map[string]bool, map[string]float64, map[string]int:
			continue
		default:
			if _, ok := v.(model.AbstractUpdateOperation); ok {
				continue
			}
			return fmt.Errorf("invalid type: %T", v)
		}
	}
	return nil
}

func BuildNewRepositoryOptionsForTest() (*NewRepositoryOption, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &NewRepositoryOption{
		DBClient: db,
	}, nil
}
