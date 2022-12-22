package databaseTools

import (
	"fmt"
	"gorm.io/gorm"
)

// Save .
func Save(db *gorm.DB, object interface{}) (rowsAffected int64, e error) {
	changeExtension, ok := object.(IChangeExtension)
	if !ok {
		return 0, fmt.Errorf("object not extend `databaseTools.ChangeExtension`")
	}

	changes := changeExtension.GetChanges()
	if changes == nil {
		return 0, nil
	}

	db = db.Model(object).Updates(changes)
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// Find .
func Find(db *gorm.DB, where, result interface{}) (e error) {
	db = db.Where(where).Find(result)

	if db.RowsAffected == 0 {
		e = gorm.ErrRecordNotFound
	} else {
		e = db.Error
	}

	return
}

// FindByPrimaryKey .
func FindByPrimaryKey(db *gorm.DB, result interface{}) (e error) {
	e = db.Last(result).Error
	return
}

// Create .
func Create(db *gorm.DB, object interface{}) (rowsAffected int64, e error) {
	db = db.Create(object)
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// BatchCreate .
func BatchCreate(db *gorm.DB, object interface{}) (rowsAffected int64, e error) {
	db = db.CreateInBatches(object, 100)
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// CreateOrSave .
func CreateOrSave(db *gorm.DB, where interface{}, data interface{}) (rowsAffected int64, e error) {
	db = db.Assign(data).FirstOrCreate(where)
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// Count .
func Count(db *gorm.DB, where interface{}) (count int64, e error) {
	e = db.
		Model(where).
		Where(where).
		Count(&count).Error

	return
}

// DeleteByPrimaryKey
//
//	删除条件只有tag上带primaryKey的成员变量
func DeleteByPrimaryKey(db *gorm.DB, object interface{}) (rowsAffected int64, e error) {
	db = db.Delete(object)
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// Delete .
func Delete(db *gorm.DB, primaryKey, where interface{}) (rowsAffected int64, e error) {
	db = db.Where(where).Delete(primaryKey)
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// BatchDeleteStringByPrimaryKey
//
//	删除条件只有tag上带primaryKey的成员变量
func BatchDeleteStringByPrimaryKey(db *gorm.DB, primaryKey interface{}, idList []string) (rowsAffected int64, e error) {
	db = db.Where(primaryKey).Delete(primaryKey, idList)
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// BatchDeleteUint8IDByPrimaryKey
//
//	删除条件只有tag上带primaryKey的成员变量
func BatchDeleteUint8IDByPrimaryKey(db *gorm.DB, primaryKey interface{}, idList []uint8) (rowsAffected int64, e error) {
	db = db.Where(primaryKey).Delete(primaryKey, idList)
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}
