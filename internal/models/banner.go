package models

import "time"

type Banner struct {
	ID        int       `json:"banner_id" db:"id"`
	Tags      []int32   `json:"tag_ids" db:"-"`
	Feature   int32     `json:"feature_id" db:"-"`
	Content   string    `json:"content" db:"content"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type BannerTags struct {
	ID        int       `json:"banner_id" db:"id"`
	Tag       *int32    `json:"tag_id" db:"tag_id"`
	Feature   int32     `json:"feature_id" db:"feature_id"`
	Content   string    `json:"content" db:"content"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
type BannerHistoryTags struct {
	ID              int       `json:"banner_id" db:"id"`
	BannerHistoryID int       `json:"banner_history_id" db:"banner_history_id"`
	Tag             *int32    `json:"tag_id" db:"tag_id"`
	Feature         int32     `json:"feature_id" db:"feature_id"`
	Content         string    `json:"content" db:"content"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
type BannerChange struct {
	ID       *int     `json:"banner_id" db:"id"`
	Tags     *[]int32 `json:"tag_id" db:"tag_id"`
	Feature  *int32   `json:"feature_id" db:"feature_id"`
	Content  *string  `json:"content" db:"content"`
	IsActive *bool    `json:"is_active" db:"is_active"`
}
