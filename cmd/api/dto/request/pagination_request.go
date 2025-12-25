package request

import (
	_ "github.com/go-playground/validator/v10"
)

type PaginationRequest struct {
	Page        int    `url:"page" validate:"required"`
	PageSize    int    `url:"page_size" validate:"required"`
	SearchField string `url:"search_field"`
	SearchValue string `url:"search_value"`
	SearchAll   string `url:"search_all"`
	OrderBy     string `url:"order_by"`
	Sort        string `url:"sort"`
}
