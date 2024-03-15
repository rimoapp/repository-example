package repository

import (
	"context"
	"reflect"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gogo/status"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
)

type firestoreGenericRepository[T model.AbstractEntity] struct {
	client         *firestore.Client
	collectionPath string // objectによっては動的にcollectionPathを変更しないときもあるので、その場合はこっちを使う
}

func newFirestoreGenericRepository[T model.AbstractEntity](client *firestore.Client, collectionPath string) *firestoreGenericRepository[T] {
	return &firestoreGenericRepository[T]{
		client:         client,
		collectionPath: collectionPath,
	}
}

var firestoreErrNilDocRef = errors.New("firestore: nil DocumentRef")

func (r *firestoreGenericRepository[T]) Get(ctx context.Context, id string) (T, error) {
	return r.get(ctx, id, r.collectionPath)
}

func (r *firestoreGenericRepository[T]) get(ctx context.Context, id, collectionPath string) (T, error) {
	var object T
	if collectionPath == "" {
		return object, errors.New("collection path is empty")
	}
	docRef := r.client.Collection(collectionPath).Doc(id)
	if docRef == nil {
		return object, status.New(codes.NotFound, "doc ref is not found").Err()
	}
	doc, err := r.client.Collection(collectionPath).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound || err == firestoreErrNilDocRef {
			return object, status.New(codes.NotFound, "not found").Err()
		}
		if strings.Contains(err.Error(), "the connection is draining") {
			return r.get(ctx, id, collectionPath)
		}
		return object, errors.Wrap(err, "failed to get")
	}
	if err := doc.DataTo(&object); err != nil {
		return object, errors.Wrap(err, "failed to assign")
	}
	object.SetID(id)
	return object, nil
}

func (r *firestoreGenericRepository[T]) Delete(ctx context.Context, id string) error {
	return r.delete(ctx, id, r.collectionPath)
}

func (r *firestoreGenericRepository[T]) delete(ctx context.Context, id, collectionPath string) error {
	if collectionPath == "" {
		return errors.New("collection path is empty")
	}
	_, err := r.client.Collection(collectionPath).Doc(id).Delete(ctx)
	return err
}
func (r *firestoreGenericRepository[T]) Create(ctx context.Context, obj T) (string, error) {
	obj.BeforeCreate(time.Now())
	return r.create(ctx, obj, r.collectionPath)
}
func (r *firestoreGenericRepository[T]) create(ctx context.Context, obj T, collectionPath string) (string, error) {
	if collectionPath == "" {
		return "", errors.New("collection path is empty")
	}
	ref, _, err := r.client.Collection(collectionPath).Add(ctx, obj)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

func (r *firestoreGenericRepository[T]) Update(ctx context.Context, id string, keyValues map[string]interface{}) error {
	return r.update(ctx, id, r.collectionPath, keyValues)
}

func (r *firestoreGenericRepository[T]) update(ctx context.Context, id, collectionPath string, keyValues map[string]interface{}) error {
	if err := validateKeyValues(keyValues); err != nil {
		return err
	}
	updates := []firestore.Update{}
	for key, value := range keyValues {
		if operation, ok := value.(model.AbstractUpdateOperation); ok {
			updates = append(updates, firestore.Update{Path: key, Value: operation.FirestoreUpdateValue()})
			continue
		}
		updates = append(updates, firestore.Update{Path: key, Value: value})
	}
	if collectionPath == "" {
		return errors.New("collection path is empty")
	}

	var object T
	updates = validateUpdateFields(updates, object)
	_, err := r.client.Collection(collectionPath).Doc(id).Update(ctx, updates)
	return err
}

func (r *firestoreGenericRepository[T]) list(ctx context.Context, query firestore.Query) ([]T, error) {
	iter := query.Documents(ctx)
	objects := []T{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to iterate")
		}

		var object T
		if err := doc.DataTo(&object); err != nil {
			return nil, errors.Wrap(err, "failed to assign")
		}
		object.SetID(doc.Ref.ID)
		objects = append(objects, object)
	}
	return objects, nil
}

func (r *firestoreGenericRepository[T]) Set(ctx context.Context, id string, obj T) error {
	return r.set(ctx, id, r.collectionPath, obj)
}

func (r *firestoreGenericRepository[T]) set(ctx context.Context, id, collectionPath string, obj T) error {
	if collectionPath == "" {
		return errors.New("collection path is empty")
	}
	if id == "" {
		id = obj.GetID()
	}
	if id == "" {
		return errors.New("id is empty")
	}
	_, err := r.client.Collection(collectionPath).Doc(id).Set(ctx, obj)
	return err
}

func listValidFields(originalType interface{}) map[string]bool {
	return listValidFieldsByType(reflect.TypeOf(originalType))
}

func listValidFieldsByType(modelType reflect.Type) map[string]bool {
	validFields := map[string]bool{}
	// ポインタ型の場合、Elem()で実際の型を取得
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	// 型が構造体でない場合は空のマップを返す
	if modelType.Kind() != reflect.Struct {
		return validFields
	}
	for i := 0; i < modelType.NumField(); i++ {
		fieldType := modelType.Field(i)
		if fieldType.Anonymous {
			embedded := listValidFieldsByType(fieldType.Type)
			for k, v := range embedded {
				validFields[k] = v
			}
		}
		tag := fieldType.Tag.Get("firestore")
		if strings.HasSuffix(tag, ",omitempty") {
			tag = strings.ReplaceAll(tag, ",omitempty", "")
		}
		validFields[tag] = true
	}
	delete(validFields, "-")
	return validFields
}

func validateUpdateFields(updates []firestore.Update, originalType interface{}) []firestore.Update {
	res := []firestore.Update{}
	validFields := listValidFields(originalType)
	updatedAtExists := false
	for _, update := range updates {
		if update.Path == "updated_at" {
			updatedAtExists = true
		}
		if validFields[update.Path] {
			res = append(res, update)
		}
	}
	if validFields["updated_at"] && !updatedAtExists {
		res = append(res, firestore.Update{Path: "updated_at", Value: time.Now()})
	}
	return res
}
