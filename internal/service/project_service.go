package service

import (
	"context"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/internal/store"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
	"go.uber.org/zap"
)

type ProjectService struct {
	logger *zap.Logger
	store  store.Storage
}

func NewProjectService(store store.Storage, logger *zap.Logger) *ProjectService {
	return &ProjectService{
		logger: logger,
		store:  store,
	}
}

func (s *ProjectService) CheckSlug(ctx context.Context, req request.CheckSlugRequest) (bool, error) {
	slugExist, err := s.store.IProject.CheckSlug(ctx, req)

	if err != nil {
		return false, err
	}

	return slugExist, nil
}

func (s *ProjectService) AddProject(ctx context.Context, userId uint, req request.AddProjectRequest) (entity.Project, error) {
	project, err := s.store.IProject.AddProject(ctx, userId, req)

	if err != nil {
		return entity.Project{}, err
	}

	return project, nil
}

func (s *ProjectService) GetProject(ctx context.Context, userId uint, req request.PaginationRequest) (utils.PaginateResult[entity.Project], error) {
	result, err := s.store.IProject.GetProject(ctx, userId, req)

	if err != nil {
		return utils.PaginateResult[entity.Project]{}, err
	}

	return result, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, projectId uint) error {
	err := s.store.IProject.DeleteProject(ctx, projectId)

	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, projectId uint, req request.AddProjectRequest) (entity.Project, error) {
	project, err := s.store.IProject.UpdateProject(ctx, projectId, req)

	if err != nil {
		return entity.Project{}, err
	}

	return project, nil
}
