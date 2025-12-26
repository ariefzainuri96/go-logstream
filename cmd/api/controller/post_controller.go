package controller

import (
	"encoding/json"
	"net/http"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/response"
	"github.com/ariefzainuri96/go-logstream/cmd/api/utils"
)

// @Summary      Add Post
// @Description  Add new Post
// @Tags         post
// @Accept       json
// @Produce      json
// @Param        request		body	  request.AddPostRequest	true "Add Post request"
// @security 	 ApiKeyAuth
// @Success      200  			{object}  response.PostResponse
// @Failure      400  			{object}  response.BaseResponse
// @Failure      404  			{object}  response.BaseResponse
// @Router       /posts/		[post]
func (app *Application) addPost(w http.ResponseWriter, r *http.Request) {
	var data request.AddPostRequest

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

	post, err := app.Service.IPost.CreatePost(r.Context(), data)

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, response.PostResponse{
		BaseResponse: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Success add post",
		},
		Post: post,
	})
}

// @Summary      Get Post
// @Description  Get All Post
// @Tags         post
// @Accept       json
// @Produce      json
// @Param        request		query	  request.GetPostRequest 	true "Get Post request"
// @security 	 ApiKeyAuth
// @Success      200  			{object}  response.PostsResponse
// @Failure      400  			{object}  response.BaseResponse
// @Failure      404  			{object}  response.BaseResponse
// @Router       /posts/		[get]
func (app *Application) getPost(w http.ResponseWriter, r *http.Request) {
	var data request.GetPostRequest

	err := decoder.Decode(&data, r.URL.Query())

	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	result, err := app.Service.IPost.GetPost(r.Context(), data)

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, response.PostsResponse{
		BaseResponse: response.BaseResponse{
			Message: "Success",
			Status:  http.StatusOK,
		},
		Posts:      result.Data,
		Pagination: result.Pagination,
	})
}

func (app *Application) PostController() *http.ServeMux {
	productRouter := http.NewServeMux()

	productRouter.HandleFunc("POST /", app.addPost)
	productRouter.HandleFunc("GET /", app.getPost)

	// Catch-all route for undefined paths
	productRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 page not found", http.StatusNotFound)
	})

	return productRouter
}
