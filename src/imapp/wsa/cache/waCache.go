package cache

import (
	"encoding/hex"
	"github.com/hyahm/golog"
)

func (c *WACache) SaveServerStaticKey(static []byte) {
	buf := hex.EncodeToString(static)
	ok := c.Cache.Set(WAKServerStaticKey, buf, 0)
	if !ok {
		golog.Warnf("save server static:%v failed.", buf)
	}
}

func (c *WACache) GetServerStaticKey() []byte {
	if val, ok := c.Cache.Get(WAKServerStaticKey); val == nil || !ok {
		return nil
	} else if val != nil {
		buff := val.(string)
		sk, _ := hex.DecodeString(buff)
		return sk
	}
	return nil
}
