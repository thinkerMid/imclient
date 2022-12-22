package databaseTools

// IChangeExtension .
type IChangeExtension interface {
	GetChanges() map[string]interface{}
}

// ChangeExtension .
type ChangeExtension struct {
	changes map[string]interface{}
}

// Update .
func (c *ChangeExtension) Update(name string, value interface{}) {
	if c.changes == nil {
		c.changes = make(map[string]interface{})
	}
	c.changes[name] = value
}

// GetChanges .
func (c *ChangeExtension) GetChanges() map[string]interface{} {
	if c.changes == nil {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range c.changes {
		result[k] = v
	}

	c.changes = nil
	return result
}
