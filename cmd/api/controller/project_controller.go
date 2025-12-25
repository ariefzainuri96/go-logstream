package controller

import (
	"encoding/json"
	"net/http"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/response"
	"github.com/ariefzainuri96/go-logstream/cmd/api/middleware"
	"github.com/ariefzainuri96/go-logstream/cmd/api/utils"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// @Summary      Add Project
// @Description  Add new Project
// @Tags         project
// @Accept       json
// @Produce      json
// @Param        request		body	  request.AddProjectRequest	true "Add Project request"
// @security 	 ApiKeyAuth
// @Success      200  			{object}  response.ProjectResponse
// @Failure      400  			{object}  response.BaseResponse
// @Failure      404  			{object}  response.BaseResponse
// @Router       /projects/		[post]
func (app *Application) addProject(w http.ResponseWriter, r *http.Request) {
	var data request.AddProjectRequest

	user, ok := middleware.GetUserFromContext(r)

	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized, please re login!")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	defer r.Body.Close()

	err = app.Validator.Struct(data)

	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	project, err := app.Service.IProject.AddProject(r.Context(), user["user_id"].(uint), data)

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, response.ProjectResponse{
		BaseResponse: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Success add project",
		},
		Project: project,
	})
}

// @Summary      Get Project
// @Description  Get All Project
// @Tags         project
// @Accept       json
// @Produce      json
// @Param        request		query	  request.PaginationRequest	true "Get Project request"
// @security 	 ApiKeyAuth
// @Success      200  			{object}  response.ProjectsResponse
// @Failure      400  			{object}  response.BaseResponse
// @Failure      404  			{object}  response.BaseResponse
// @Router       /projects/		[get]
func (app *Application) getProject(w http.ResponseWriter, r *http.Request) {
	var data request.PaginationRequest

	err := decoder.Decode(&data, r.URL.Query())

	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

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

// @Summary      Delete Project
// @Description  Delete Project
// @Tags         project
// @Produce      json
// @Param        id   				path      int  true  "Project ID"
// @security 	 ApiKeyAuth
// @Success      200  				{object}  response.BaseResponse
// @Failure      400  				{object}  response.BaseResponse
// @Failure      404  				{object}  response.BaseResponse
// @Router       /projects/{id}		[delete]
func (app *Application) deleteProject(w http.ResponseWriter, r *http.Request) {
	// id, err := strconv.Atoi(r.PathValue("id"))

	// if err != nil {
	// 	utils.RespondError(w, http.StatusBadRequest, "Invalid id")
	// 	return
	// }

	// err = app.store.IProduct.DeleteProduct(r.Context(), uint(id))

	// if err != nil {
	// 	utils.RespondError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// utils.WriteJSON(w, http.StatusOK, response.BaseResponse{
	// 	Status:  http.StatusOK,
	// 	Message: "Success delete product",
	// })
}

// @Summary      Update Project
// @Description  Update Project
// @Tags         project
// @Accept       json
// @Produce      json
// @Param 		 id					path      int  true  "Project ID"
// @Param        request			body	  request.AddProjectRequest	true "Add Project request"
// @security 	 ApiKeyAuth
// @Success      200  				{object}  response.ProjectResponse
// @Failure      400  				{object}  response.BaseResponse
// @Failure      404  				{object}  response.BaseResponse
// @Router       /projects/{id}		[put]
func (app *Application) updateProject(w http.ResponseWriter, r *http.Request) {
	// productID, err := strconv.Atoi(r.PathValue("id"))

	// if err != nil {
	// 	utils.RespondError(w, http.StatusBadRequest, "Invalid id")
	// 	return
	// }

	// // Decode request body into a map
	// var updateData map[string]any
	// err = json.NewDecoder(r.Body).Decode(&updateData)
	// if err != nil {
	// 	utils.RespondError(w, http.StatusBadRequest, "Invalid request")
	// 	return
	// }
	// defer r.Body.Close()

	// // Ensure there's data to update
	// if len(updateData) == 0 {
	// 	http.Error(w, "No fields to update", http.StatusBadRequest)
	// 	return
	// }

	// product, err := app.store.IProduct.PatchProduct(r.Context(), uint(productID), updateData)

	// if err != nil {
	// 	utils.RespondError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// utils.WriteJSON(w, http.StatusOK, response.ProductResponse{
	// 	BaseResponse: response.BaseResponse{
	// 		Status:  http.StatusOK,
	// 		Message: "Success patch product",
	// 	},
	// 	Product: product,
	// })
}

func (app *Application) ProjectController() *http.ServeMux {
	productRouter := http.NewServeMux()

	productRouter.HandleFunc("POST /", app.addProject)
	productRouter.HandleFunc("GET /", app.getProject)
	productRouter.HandleFunc("DELETE /{id}", app.deleteProject)
	productRouter.HandleFunc("PUT /{id}", app.updateProject)

	// Catch-all route for undefined paths
	productRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 page not found", http.StatusNotFound)
	})

	return productRouter
}
