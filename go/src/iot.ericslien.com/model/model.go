package model

import (
	"errors"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"iot.ericslien.com/database"
)

type CrudBase interface {
	// @todo: We should be able to have this be a more specific/accurate interface
	GetId() uint64
	StrToId(string) uint64
	IdToStr(uint64) string

	Search(params map[string]interface{}, limit int, offset int) ([]CrudBase, error)
	Load(id uint64) error
	Create() error
	Update() error
	Delete() error

	// Should be handled elsewhere? Is a model the right place?
	// Access(perm string, user *UserTokenClaims) (bool, error)
}

type NewModel func() CrudBase

//Common Crud implementation helpers using gorm.

func Load(model CrudBase, id uint64) error {
	db, _ := database.Open()
	return db.First(model, id).Error
}

func Search(objs interface{}, params map[string]interface{}, limit int, offset int) error {
	db, _ := database.Open()
	if limit > 0 {
		db = db.Limit(limit)
	}

	if offset > 0 {
		if limit == 0 {
			return errors.New("Cannot have an offset without a limit.")
		}
		db = db.Offset(offset)
	}

	err := db.Find(objs, params).Error
	return err
}

func Create(model CrudBase) error {
	db, _ := database.Open()
	return db.Create(model).Error
}

func Update(model CrudBase) error {
	db, _ := database.Open()
	return db.Model(model).Updates(model).Error
}

func Delete(model CrudBase) error {
	db, _ := database.Open()
	if model.GetId() != 0 {
		return db.Delete(model, model).Error
	} else {
		return errors.New("Cannot delete object if no ID is set")
	}
}

func FilterSearchParams(obj interface{}, params map[string]interface{}) map[string]interface{} {
	s := structs.New(obj)

	db_field_names := make(map[string]bool, 0)
	for _, f := range s.Fields() {
		db_field_names[gorm.ToDBName(f.Name())] = true
	}

	for k, _ := range params {
		if !db_field_names[k] {
			delete(params, k)
		}
	}
	return params
}
