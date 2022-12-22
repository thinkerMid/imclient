package businessDB

import databaseTools "ws/framework/plugin/database/database_tools"

// BusinessProfile .
type BusinessProfile struct {
	databaseTools.ChangeExtension

	JID              string `gorm:"column:jid;primaryKey"`
	ProfileTag       string `gorm:"column:profile_tag"`
	CatalogSessionID string `gorm:"column:catalog_session_id"`
	VNameSerial      int64  `gorm:"column:vname_serial"`
}

// TableName .
func (c *BusinessProfile) TableName() string {
	return "business"
}

// UpdateProfileTag .
func (c *BusinessProfile) UpdateProfileTag(v string) {
	if v == c.ProfileTag {
		return
	}

	c.ProfileTag = v
	c.Update("profile_tag", v)
}

// UpdateCatalogSessionID .
func (c *BusinessProfile) UpdateCatalogSessionID(v string) {
	if v == c.CatalogSessionID {
		return
	}

	c.CatalogSessionID = v
	c.Update("catalog_session_id", v)
}

// UpdateVNameSerial .
func (c *BusinessProfile) UpdateVNameSerial(v int64) {
	if v == c.VNameSerial {
		return
	}

	c.VNameSerial = v
	c.Update("vname_serial", v)
}
