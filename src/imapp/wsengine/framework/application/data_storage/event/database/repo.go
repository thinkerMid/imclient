package eventDB

import "gorm.io/gorm"

// FindLast .
func FindLast(db *gorm.DB, jid string, serialNumber int64, channelID byte) (results []EventBuffer, e error) {
	e = db.
		Where("jid", jid).
		Where("channel_id", channelID).
		Where("serial_number < ?", serialNumber).
		Order("auto_increment_id asc").
		Find(&results).Error

	return
}

// FindAll .
func FindAll(db *gorm.DB, jid string, serialNumber int64, channelID byte) (results []EventBuffer, e error) {
	e = db.
		Where("jid", jid).
		Where("channel_id", channelID).
		Where("serial_number <= ?", serialNumber).
		Order("auto_increment_id asc").
		Find(&results).Error

	return
}

// DeleteLast .
func DeleteLast(db *gorm.DB, jid string, serialNumber int64, channelID byte) (e error) {
	object := EventInfo{JID: jid}

	return db.
		Where("channel_id", channelID).
		Where("serial_number < ?", serialNumber).
		Delete(&object).Error
}

// Delete .
func Delete(db *gorm.DB, jid string, serialNumber int64, channelID byte) (e error) {
	object := EventInfo{JID: jid}

	return db.
		Where("channel_id", channelID).
		Where("serial_number <= ?", serialNumber).
		Delete(&object).Error
}

// DeleteByJID .
func DeleteByJID(db *gorm.DB, jid string) (e error) {
	object := EventInfo{JID: jid}

	return db.Delete(&object).Error
}

// Count .
func Count(db *gorm.DB, jid string, serialNumber int64, channelID byte) (count int64, e error) {
	object := EventInfo{}

	e = db.
		Model(&object).
		Where("jid", jid).
		Where("channel_id", channelID).
		Where("serial_number <= ?", serialNumber).
		Count(&count).Error
	return
}
