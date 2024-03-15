package model

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// AbstractDataEntity used in repository
type AbstractDataEntity interface {
	GetID() string
	SetID(id string)
}

type AbstractAssociatedEntity interface {
	AbstractDataEntity
	IsAuthorized(userID string) bool
	SetCreatorID(userID string)
}

type BaseDataEntity struct {
	ID        string    `json:"id" firestore:"-"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`

	// for gorm
	UID       int            `json:"-" firestore:"-" gorm:"primarykey"`
	DeletedAt gorm.DeletedAt `json:"-" firestore:"-" gorm:"index"`
}

func (a *BaseDataEntity) SetID(id string) {
	a.ID = id
	// for gorm
	uid, err := strconv.Atoi(id)
	if err == nil {
		a.UID = uid
	}
}
func (a *BaseDataEntity) GetID() string {
	// for gorm
	if a.ID == "" && a.UID != 0 {
		return fmt.Sprintf("%d", a.UID)
	}
	return a.ID
}

var _ AbstractDataEntity = (*BaseDataEntity)(nil)

type BaseAssociatedEntity struct {
	BaseDataEntity
	UserID string `json:"user_id" firestore:"user_id"`
}

func (a *BaseAssociatedEntity) IsAuthorized(userID string) bool {
	return a.UserID == userID
}

func (a *BaseAssociatedEntity) SetCreatorID(userID string) {
	a.UserID = userID
}

// Update Operation for repository

type AbstractUpdateOperation interface {
	FirestoreUpdateValue() interface{}
	ReflectValue(rv reflect.Value) (*reflect.Value, error)
}

type incrementOperation struct {
	Value interface{}
}

func (i *incrementOperation) FirestoreUpdateValue() interface{} {
	return firestore.Increment(i.Value)
}
func (i *incrementOperation) ReflectValue(current reflect.Value) (*reflect.Value, error) {
	currentType := current.Type()
	iValue := reflect.ValueOf(i.Value)
	if !iValue.Type().ConvertibleTo(currentType) {
		return nil, errors.New("type mismatch")
	}
	// diffはi.Valueをcurrentの型に変換したもの
	diff := iValue.Convert(currentType)

	// currentの値とdiffを足し合わせます
	var result reflect.Value
	switch current.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sum := current.Int() + diff.Int()
		result = reflect.New(currentType).Elem()
		result.SetInt(sum)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sum := current.Uint() + diff.Uint()
		result = reflect.New(currentType).Elem()
		result.SetUint(sum)
	case reflect.Float32, reflect.Float64:
		sum := current.Float() + diff.Float()
		result = reflect.New(currentType).Elem()
		result.SetFloat(sum)
	default:
		return nil, errors.New("unsupported type for increment operation")
	}

	// 結果を元の型に適合させます
	if result.CanConvert(currentType) {
		convertedResult := result.Convert(currentType)
		return &convertedResult, nil
	}
	return nil, errors.New("failed to convert the result back to the original type")
}
func IncrementOperation(value interface{}) *incrementOperation {
	return &incrementOperation{Value: value}
}

type arrayUnionOperation struct {
	Values []interface{}
}

func (a *arrayUnionOperation) FirestoreUpdateValue() interface{} {
	return firestore.ArrayUnion(a.Values...)
}
func (a *arrayUnionOperation) ReflectValue(current reflect.Value) (*reflect.Value, error) {
	if current.Kind() != reflect.Slice && current.Kind() != reflect.Array {
		return nil, errors.New("current is not a slice or array")
	}
	elemType := current.Type().Elem()
	for _, i := range a.Values {
		value := reflect.ValueOf(i)
		if value.Type().ConvertibleTo(elemType) {
			convertedValue := value.Convert(elemType)
			current = reflect.Append(current, convertedValue)
		} else {
			return nil, fmt.Errorf("type mismatch for array union operation: %v, %v(%v)", elemType, value, value.Type())
		}
	}

	return &current, nil
}

func ArrayUnionOperation(values ...interface{}) *arrayUnionOperation {
	return &arrayUnionOperation{Values: values}
}

type arrayRemoveOperation struct {
	Values []interface{}
}

func (a *arrayRemoveOperation) FirestoreUpdateValue() interface{} {
	return firestore.ArrayRemove(a.Values...)
}
func (a *arrayRemoveOperation) ReflectValue(current reflect.Value) (*reflect.Value, error) {
	if current.Kind() != reflect.Slice && current.Kind() != reflect.Array {
		return nil, errors.New("current is not a slice or array")
	}
	elemType := current.Type().Elem()
	newSlice := reflect.MakeSlice(current.Type(), 0, current.Len())

	for j := 0; j < current.Len(); j++ {
		currentElem := current.Index(j)
		shouldRemove := false
		for _, v := range a.Values {
			value := reflect.ValueOf(v)
			if value.Type().ConvertibleTo(elemType) {
				convertedValue := value.Convert(elemType)
				if reflect.DeepEqual(currentElem.Interface(), convertedValue.Interface()) {
					shouldRemove = true
					break
				}
			} else {
				return nil, fmt.Errorf("type mismatch for array remove operation: expected %v, got %v(%v)", elemType, value, value.Type())
			}
		}
		if !shouldRemove {
			newSlice = reflect.Append(newSlice, currentElem)
		}
	}

	return &newSlice, nil
}

func ArrayRemoveOperation(values ...interface{}) *arrayRemoveOperation {
	return &arrayRemoveOperation{Values: values}
}

var _ AbstractUpdateOperation = (*incrementOperation)(nil)
var _ AbstractUpdateOperation = (*arrayUnionOperation)(nil)
var _ AbstractUpdateOperation = (*arrayRemoveOperation)(nil)
