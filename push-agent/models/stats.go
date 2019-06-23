package models

type (
	ChannelStatsDetailed struct {
		Channel           string `json:"channel"`
		PublishedMessages int64  `json:"published_messages"`
		StoredMessages    int64  `json:"stored_messages"`
		Subscribers       int64  `json:"subscribers"`
	}

	GlobalStatsDetailed struct {
		Hostname         string                 `json:"hostname"`
		Time             string                 `json:"time"`
		Channels         int64                  `json:"channels"`
		WildcardChannels int64                  `json:"wildcard_channels"`
		Uptime           int64                  `json:"uptime"`
		Infos            []ChannelStatsDetailed `json:"infos"`
	}

	GlobalStatsSummarized struct {
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
	}

	GlobalStats struct {
		Detailed *GlobalStatsDetailed `json:"detailed"`
		Summarized *GlobalStatsSummarized `json:"summarized"`
	}
)
