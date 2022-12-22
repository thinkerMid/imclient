package senderKeyDB

import (
	"gorm.io/gorm"
	_ "ws/framework/plugin/database/database_tools"
)

// DeleteByJID .
func DeleteByJID(db *gorm.DB, jid string) (rowsAffected int64, err error) {
	object := SenderKey{}

	db = db.Where("our_jid", jid).Delete(&object)

	err = db.Error
	rowsAffected = db.RowsAffected
	return
}
