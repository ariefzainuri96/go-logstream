package request

// @Model
type GetPostRequest struct {
	PaginationRequest
	ProjectId uint `url:"project_id"`
}