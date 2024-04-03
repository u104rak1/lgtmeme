package dto

import "github.com/google/uuid"

type GetImagesQuery struct {
	Page             int      `query:"page" validate:"min=0"  default:"0"`
	Keyword          string   `query:"keyword" validate:"max=50"  default:""`
	Sort             string   `query:"sort" validate:"sort" default:"latest"`
	FavoriteImageIDs []string `query:"favorite_image_ids" validate:"uuidSlice"`
	AuthCheck        bool     `query:"auth_check"`
}

type GetImagesImages struct {
	ID  uuid.UUID `json:"id"`
	URL string    `json:"url"`
}

type GetImagesResp struct {
	Total  int               `json:"total"`
	Images []GetImagesImages `json:"images"`
}
