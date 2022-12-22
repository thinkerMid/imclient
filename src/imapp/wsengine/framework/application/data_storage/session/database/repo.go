package sessionDB

import (
	"gorm.io/gorm"
	_ "ws/framework/plugin/database/database_tools"
)

// DeleteByJID .
func DeleteByJID(db *gorm.DB, jid string) (rowsAffected int64, err error) {
	object := Session{}

	db = db.Where("ourJid", jid).Delete(&object)

	err = db.Error
	rowsAffected = db.RowsAffected
	return
}
