package main

import (
	"math/rand"
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

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

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
type GetPageMetaRequest struct {
	Url string `json:"url"`
}

func NewAccount(firstName string, lastName string) *Account {
	return &Account{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(1000000)),
		CreatedAt: time.Now().UTC(),
	}
}
