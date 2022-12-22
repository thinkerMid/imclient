package accountDB

import (
	"fmt"
	"gorm.io/gorm"
	"ws/framework/application/data_storage/account/constant"
)

// DeleteByJID .
func DeleteByJID(db *gorm.DB, id string) (e error) {
	object := Account{JID: id}

	e = db.Delete(&object).Error
	return
}

// FindUnavailableAccounts .
func FindUnavailableAccounts(db *gorm.DB, timeLimit int64) (results []AccountJID, e error) {
	model := AccountJID{}

	e = db.Model(&model).
		Where("(status between 401 and 503").
		Or(fmt.Sprintf("status = %v)", accountServiceConstant.Unregistered)).
		Where(fmt.Sprintf("logout_time between 0 and %v", timeLimit)).
		Limit(500).
		Find(&results).Error
	return
}
