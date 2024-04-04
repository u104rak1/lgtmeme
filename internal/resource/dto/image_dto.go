package dto

import "github.com/google/uuid"

type PostImageReqBody struct {
	Base64image string `json:"base64image" validate:"required,imageSize,base64image"`
	Keyword     string `json:"keyword" validate:"max=50"`
}

type PostImageResp struct {
	ImageURL string `json:"imageUrl"`
}

type GetImagesQuery struct {
	Page             int    `query:"page" validate:"min=0"  default:"0"`
	Keyword          string `query:"keyword" validate:"max=50" default:""`
	Sort             string `query:"sort" validate:"sort" default:"latest"`
	FavoriteImageIDs string `query:"favorite_image_ids" validate:"uuidStrings"`
	AuthCheck        bool   `query:"auth_check"`
}

type GetImagesImages struct {
	ID  uuid.UUID `json:"id"`
	URL string    `json:"url"`
}

type GetImagesResp struct {
	Total  int               `json:"total"`
	Images []GetImagesImages `json:"images"`
}
