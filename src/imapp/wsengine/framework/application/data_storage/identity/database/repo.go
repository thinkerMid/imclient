package identityDB

import (
	"gorm.io/gorm"
)

// DeleteByJID .
func DeleteByJID(db *gorm.DB, jid string) (rowsAffected int64, err error) {
	object := Identity{}

	db = db.Where("ourJid", jid).Delete(&object)

	err = db.Error
	rowsAffected = db.RowsAffected
	return
}
