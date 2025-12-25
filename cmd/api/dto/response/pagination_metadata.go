package response

type PaginationMetadata struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	TotalPage int   `json:"total_page"`
	TotalData int64 `json:"total_data"`
}
