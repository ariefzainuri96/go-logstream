package interfaces

import (
	"context"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
)

type IProject interface {
	CheckSlug(context.Context, request.CheckSlugRequest) (bool, error)
	AddProject(context.Context, uint, request.AddProjectRequest) (entity.Project, error)
	GetProject(context.Context, uint, request.PaginationRequest) (utils.PaginateResult[entity.Project], error)
	DeleteProject(context.Context, uint) error
	UpdateProject(context.Context, uint, request.AddProjectRequest) (entity.Project, error)
}
