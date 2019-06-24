package models

type (
	ChannelStats struct {
		Channel           string `json:"channel"`
		PublishedMessages int64  `json:"published_messages"`
		StoredMessages    int64  `json:"stored_messages"`
		Subscribers       int64  `json:"subscribers"`
		Hostname          string `json:"hostname"` // set by us, doesn't come from push-stream
		Agent             string `json:"agent"`    // set by us, doesn't come from push-stream
	}
)
