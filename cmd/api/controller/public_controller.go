package controller

import "net/http"

// @Summary      Get All Project
// @Description  Get All Project
// @Tags         public
// @Accept       json
// @Produce      json
// @Param        request					query	  request.PaginationRequest	true "Get Project request"
// @security 	 ApiKeyAuth
// @Success      200  						{object}  response.ProjectsResponse
// @Failure      400  						{object}  response.BaseResponse
// @Failure      404  						{object}  response.BaseResponse
// @Router       /public/projects/{projectId}		[get]
func (app *Application) getPublicProject(w http.ResponseWriter, r *http.Request) {
	// var data request.PaginationRequest

	// err := decoder.Decode(&data, r.URL.Query())

	// if err != nil {
	// 	utils.RespondError(w, http.StatusBadRequest, "Invalid request")
	// 	return
	// }

	// product, err := app.store.IProduct.GetProduct(r.Context(), data)

	// if err != nil {
	// 	utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
	// 	return
	// }

	// utils.WriteJSON(w, http.StatusOK, response.ProductsResponse{
	// 	BaseResponse: response.BaseResponse{
	// 		Message: "Success",
	// 		Status:  http.StatusOK,
	// 	},
	// 	Products:   product.Data,
	// 	Pagination: product.Pagination,
	// })
}

func (app *Application) PublicController() *http.ServeMux {
	productRouter := http.NewServeMux()

	productRouter.HandleFunc("GET /projects/{projectId}", app.getPublicProject)

	// Catch-all route for undefined paths
	productRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 page not found", http.StatusNotFound)
	})

	return productRouter
}
