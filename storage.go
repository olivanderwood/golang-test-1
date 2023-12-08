package main

import (
	"database/sql"
	"fmt"
	// "os"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetRandomApiKey() (*ApiKey, error)
	GetPageMetaByUrl(string) (*PageMeta, error)
	CreatePageMeta(*PageMeta) error
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) Init() error {
	err := s.createApiKeyTable()
	if err != nil {
		return err
	}

	return s.createPageMetaTable()
}

func (s *PostgresStore) createPageMetaTable() error {
	query := `create table if not exists page_meta (
		id serial primary key,
		url varchar(500),
		type varchar(100),
		version varchar(100),
		title text,
		author varchar(50),
		provider_name varchar(100),
		thumbnail_url varchar(500),
		thumbnail_width integer,
		thumbnail_height integer,
		html text,
		cache_age integer,
		data_iframely_url boolean,
		youtube_video_id text,
		description text
	)`

	_, err := s.db.Exec(query)
	return err
}
func (s *PostgresStore) createApiKeyTable() error {
	query := `create table if not exists api_key (
		id serial primary key,
		key varchar(100),
		usage_count integer
	)`
// 	insertQuery := `insert into api_key (id, key, usage_count) values (1, '1d967a0a2bdbd3e0b72b4f', 0);
// insert into api_key (id, key, usage_count) values (2, 'e9857f009023ca6eba88df', 0);
// insert into api_key (id, key, usage_count) values (3, '95a07c09131158b0c0b377', 0);
// insert into api_key (id, key, usage_count) values (4, 'd61aa68a1fb06483d84901', 0);
// insert into api_key (id, key, usage_count) values (5, 'e2eac4d6d9e9ca9d85b7d9', 0);`
	_, err := s.db.Exec(query)
	// _, errIns := s.db.Exec(insertQuery)
	// if errIns != nil {
	// 	return errIns
	// }
	return err
}

func NewPostgresStore() (*PostgresStore, error) {
	// connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) GetRandomApiKey() (*ApiKey, error) {
	row, err := s.db.Query(`select * from api_key order by random() limit 1`)
	apiKey, apiKeyErr := scanIntoApiKey(row)
	_, err = s.db.Exec(`update api_key set usage_count=$1
							where id=$2`, apiKey.UsageCount+1, apiKey.ID)
	if err != nil {
		return nil, err
	}
	return apiKey, apiKeyErr
}

func scanIntoApiKey(rows *sql.Rows) (*ApiKey, error) {

	for rows.Next() {
		apiKey := new(ApiKey)
		err := rows.Scan(
			&apiKey.ID,
			&apiKey.Key,
			&apiKey.UsageCount,
		)
		return apiKey, err
	}
	return nil, nil
}

func scanIntoPageMeta(rows *sql.Rows) (*PageMeta, error) {

	for rows.Next() {
		pageMeta := new(PageMeta)

		err := rows.Scan(
			&pageMeta.ID,
			&pageMeta.Url,
			&pageMeta.Type,
			&pageMeta.Version,
			&pageMeta.Title,
			&pageMeta.Author,
			&pageMeta.ProviderName,
			&pageMeta.ThumbnailUrl,
			&pageMeta.ThumbnailWidth,
			&pageMeta.ThumbnailHeight,
			&pageMeta.Html,
			&pageMeta.CacheAge,
			&pageMeta.DataIframelyUrl,
			&pageMeta.YoutubeVideoId,
			&pageMeta.Description,
		)
		return pageMeta, err

	}
	return nil, nil
}

func (s *PostgresStore) GetPageMetaByUrl(url string) (*PageMeta, error) {
	row, err := s.db.Query(`select * from page_meta where url = $1`, url)
	if err != nil {
		return nil, err
	}
	return scanIntoPageMeta(row)

}

func (s *PostgresStore) CreatePageMeta(pageMeta *PageMeta) error {
	resp, err := s.db.Exec(`insert into page_meta (author, cache_age, data_iframely_url, description, html, type, provider_name, thumbnail_height, thumbnail_width, thumbnail_url, url, title, version, youtube_video_id)
							values
							($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		pageMeta.Author, pageMeta.CacheAge, pageMeta.DataIframelyUrl, pageMeta.Description, pageMeta.Html, pageMeta.Type, pageMeta.ProviderName, pageMeta.ThumbnailHeight, pageMeta.ThumbnailWidth, pageMeta.ThumbnailUrl, pageMeta.Url, pageMeta.Title, pageMeta.Version, pageMeta.YoutubeVideoId)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}
