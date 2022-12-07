package utils

import (
	"github.com/google/uuid"
	"strings"
)

// GenUUID4 生成uuid4
func GenUUID4() string {
	u4 := uuid.New()
	return strings.ToUpper(u4.String())
}

func ParseUUID4(content string) []byte {
	b, _ := uuid.Parse(content)
	return b[:]
}
