package groupDB

import (
	"gorm.io/gorm"
)

// DeleteByJID .
func DeleteByJID(db *gorm.DB, jid string) (rowsAffected int64, err error) {
	object := Group{}

	db = db.Where("jid", jid).Delete(&object)

	err = db.Error
	rowsAffected = db.RowsAffected
	return
}
