package types

type Presence string

const (
	PresenceAvailable   Presence = "available"
	PresenceUnavailable Presence = "unavailable"
)

type ChatPresence string

const (
	ChatPresenceComposing ChatPresence = "composing"
	ChatPresenceRecording ChatPresence = "recording"
	ChatPresencePaused    ChatPresence = "paused"
)
