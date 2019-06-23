package models

type (
	// as returned by push-stream stats
	ChannelInfo struct {
		Channel string `json:"channel"`
		PublishedMessages int64 `json:"published_messages"`
		StoredMessages int64 `json:"stored_messages"`
		Subscribers int64 `json:"subscribers"`
	}
)

