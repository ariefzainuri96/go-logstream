package interfaces

import (
	"context"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
)

type IPost interface {
	CreatePost(context.Context, request.AddPostRequest) (entity.Post, error)
	GetPost(context.Context, request.GetPostRequest) (utils.PaginateResult[entity.Post], error)
}
