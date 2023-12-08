package main



type PageMeta struct {
	ID              int    `json:"id"`
	Url             string `json:"url"`
	Type            string `json:"type"`
	Version         string `json:"version"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	ProviderName    string `json:"procvider_name"`
	ThumbnailUrl    string `json:"thumbnail_url"`
	ThumbnailWidth  int64  `json:"thumbnail_width"`
	ThumbnailHeight int64  `json:"thumbnail_height"`
	Html            string `json:"html"`
	CacheAge        int64  `json:"cache_age"`
	DataIframelyUrl bool   `json:"data_iframely_url"`
	YoutubeVideoId  string `json:"youtube_video_id"`
	Description     string `json:"description"`
}

type ApiKey struct {
	ID         int    `json:"id"`
	UsageCount int    `json:"usage_count"`
	Key        string `json:"key"`
}


type GetPageMetaRequest struct {
	Url string `json:"url"`
}

