package models

type (
	GlobalStats struct {
		Hostname          string `json:"hostname"`
		Time              string `json:"time"`
		Channels          int64  `json:"channels"`
		WildcardChannels  int64  `json:"wildcard_channels"`
		PublishedMessages int64  `json:"published_messages"`
		StoredMessages    int64  `json:"stored_messages"`
		MessagesInTrash   int64  `json:"messages_in_trash"`
		ChannelsInDelete  int64  `json:"channels_in_delete"`
		ChannelsInTrash   int64  `json:"channels_in_trash"`
		Subscribers       int64  `json:"subscribers"`
		Uptime            int64  `json:"uptime"`
		Agent             string `json:"agent"`    // set by us, doesn't come from push-stream
	}

	// returned by push-stream with more fields, but we only use Hostname and Infos
	GlobalStatsDetailed struct {
		Hostname string         `json:"hostname"`
		Infos    []ChannelStats `json:"infos"`
	}
)
