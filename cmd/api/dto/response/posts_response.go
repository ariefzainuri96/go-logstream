package response

import "github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"

// @Model
type PostsResponse struct {
	BaseResponse
	Posts []entity.Post `json:"posts"`
	Pagination PaginationMetadata `json:"pagination"`
}

// @Model
type PostResponse struct {
	BaseResponse
	Post entity.Post `json:"post"`
}