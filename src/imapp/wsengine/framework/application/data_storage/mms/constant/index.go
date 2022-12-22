package mmsConstant

// Bucket .
type Bucket struct {
	Host         string // 原地址
	FallbackHost string // dns被污染或者无法解析 使用fallback地址
}

// MMSEndpoints .
type MMSEndpoints struct {
	ID      string
	Auth    string
	Buckets []Bucket
}
