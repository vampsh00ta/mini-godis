package request

type GetBannerForUser struct {
	TagID           int32 `json:"tag_id" validate:"required" schema:"tag_id"`
	FeatureID       int32 `json:"feature_id" validate:"required" schema:"feature_id"`
	UseLastRevision bool  `json:"use_last_revision" schema:"use_last_revision"`
}

type GetBanners struct {
	TagID     int32 `json:"tag_id" schema:"tag_id"`
	FeatureID int32 `json:"feature_id" schema:"feature_id"`
	Limit     int32 `json:"limit" schema:"limit"`
	Offset    int32 `json:"offset" schema:"offset"`
}
type GetBannerHistory struct {
	Limit int `json:"limit" schema:"limit"`
}
type DeleteBannerByID struct {
	ID int `json:"id"  `
}

type CreateBanner struct {
	Tags     []int32 `json:"tag_ids" validate:"required,gt=0,dive,required" `
	Feature  int32   `json:"feature_id"  validate:"required"`
	Content  string  `json:"content"  validate:"required"`
	IsActive bool    `json:"is_active" validate:"required" `
}

type ChangeBanner struct {
	Tags     *[]int32 `json:"tag_ids"  `
	Feature  *int32   `json:"feature_id"  `
	Content  *string  `json:"content"  `
	IsActive *bool    `json:"is_active"  `
}
type DeleteBannerByTagAndFeature struct {
	TagID     int32 `json:"tag_id" validate:"required" schema:"tag_id" `
	FeatureID int32 `json:"feature_id"  validate:"required" schema:"feature_id"`
}
