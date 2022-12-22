package events

//type Option interface {
//}

type MediaType int

const (
	MediaText MediaType = iota + 1
	MediaImage
	MediaVideo
	MediaVoice
	MediaVCard
	MediaDocument MediaType = 8
)

type MessageType int

const (
	MessagePrivate MessageType = 1
	MessageGroup
)
