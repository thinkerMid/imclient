package deviceListDB

import (
	"gorm.io/gorm"
)

// DeleteByJID .
func DeleteByJID(db *gorm.DB, id string) (rowsAffected int64, err error) {
	object := Device{}

	db = db.Where("our_jid", id).Delete(&object)
	err = db.Error
	rowsAffected = db.RowsAffected
	return
}
