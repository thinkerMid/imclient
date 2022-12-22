package prekeyDB

import (
	"gorm.io/gorm"
)

// FindAll .
func FindAll(db *gorm.DB, where *PreKey, limit int) (results []PreKey, e error) {
	e = db.
		Model(where).
		Where(where).
		Order("keyId desc").
		Limit(limit).
		Find(&results).Error

	return
}

// FindLast .
func FindLast(db *gorm.DB, where *PreKey) (e error) {
	e = db.
		Order("keyId desc").
		Limit(1).
		Find(where).Error

	return
}

// DeleteByKeyIDs .
func DeleteByKeyIDs(db *gorm.DB, where *PreKey, preKeyIDs []uint32) (e error) {
	e = db.
		Where("keyId in (?)", preKeyIDs).
		Delete(&where).Error
	return
}
