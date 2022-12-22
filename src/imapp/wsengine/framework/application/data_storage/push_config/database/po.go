package pushConfigDB

import (
	"time"
	"ws/framework/plugin/database/database_tools"
	"ws/framework/utils"
)

var monthIdx uint8 = 31

// PushConfig .
type PushConfig struct {
	databaseTools.ChangeExtension

	JID       string `gorm:"column:jid;primaryKey"`
	VoipToken string `gorm:"column:voip_token"`
	ApnsToken string `gorm:"column:apns_token"`
	Pkey      []byte `gorm:"column:pkey"`
}

// NewPushConfig .
func NewPushConfig(JID string) *PushConfig {
	p := PushConfig{JID: JID}
	p.newPkey()

	p.ApnsToken = utils.RandPushToken()
	p.VoipToken = utils.RandPushToken()

	return &p
}

// TableName .
func (s *PushConfig) TableName() string {
	return "push_config"
}

func (s *PushConfig) newPkey() {
	month := uint8(time.Now().Month())

	pkeyBuffer := utils.RandBytes(32)
	pkeyBuffer[monthIdx] = month
	s.Pkey = pkeyBuffer
}

// RefreshPkey
//
//	pkey有定时更新逻辑 真机是30天更新一次 目前逻辑是按每月计算 先随机32位byte后取第32位作为月份储存
func (s *PushConfig) RefreshPkey() {
	// 当前的月份
	month := uint8(time.Now().Month())

	// 没有pkey 或 pkey的月份已经不是当月的
	if len(s.Pkey) == 0 || s.Pkey[monthIdx] != month {
		s.newPkey()

		s.Update("pkey", s.Pkey)
	}
}
