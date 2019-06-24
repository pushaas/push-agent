package models

type (
	ChannelStatsDetailed struct {
		// come from push-stream
		Channel           string `json:"channel"`
		PublishedMessages int64  `json:"published_messages"`
		StoredMessages    int64  `json:"stored_messages"`
		Subscribers       int64  `json:"subscribers"`

		// TODO remove if not needed
		// set by us
		//Updated  time.Time `json:"updated"`
		Hostname string    `json:"hostname"`
	}

	// returned by push-stream with more fields, but we only use Hostname and Infos
	GlobalStatsDetailed struct {
		Hostname string                 `json:"hostname"`
		Infos    []ChannelStatsDetailed `json:"infos"`
	}

	GlobalStatsSummarized struct {
		// come from push-stream
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

		// TODO remove if not needed
		// set by us
		//Updated time.Time `json:"updated"`
	}
)
