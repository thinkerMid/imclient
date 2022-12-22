package contactDB

import (
	"gorm.io/gorm"
)

// FindList .
func FindList(db *gorm.DB, srcNumber string, searchJIDNumber []string) (result []Contact, err error) {
	db = db.Where("src_number", srcNumber).Where("dst_jid_number", searchJIDNumber).Find(&result)
	err = db.Error
	return
}

// Delete .
func Delete(db *gorm.DB, object *Contact) (e error) {
	if len(object.DstPhoneNumber) > 0 {
		db = db.Model(object).Where("src_number", object.JID).Where("dst_phone_number", object.DstPhoneNumber)
	} else {
		db = db.Model(object).Where("src_number", object.JID).Where("dst_jid_number", object.DstJIDUser)
	}

	query := Contact{}

	e = db.Delete(&query).Error
	return
}

// DeleteByJID .
func DeleteByJID(db *gorm.DB, id string) (rowsAffected int64, err error) {
	object := Contact{}

	db = db.Where("src_number", id).Delete(&object)
	err = db.Error
	rowsAffected = db.RowsAffected
	return
}
