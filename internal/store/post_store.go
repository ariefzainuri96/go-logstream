package store

import (
	"context"
	"time"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/internal/db"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostStore struct {
	gormDb *db.GormDB
	logger *zap.Logger
}

func (s *PostStore) CreatePost(ctx context.Context, req request.AddPostRequest) (entity.Post, error) {
	post := entity.Post{
		ProjectId: req.ProjectId,
		Title:     req.Title,
		Content:   req.Content,
		Category:  req.Category,
		Status:    req.Status,
	}

	result := s.gormDb.ExecWithTimeoutVal(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Create(&post)
	})

	if result.Error != nil {
		return entity.Post{}, result.Error
	}

	var project entity.Project

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.First(&project, req.ProjectId).Error
	})

	if err != nil {
		return entity.Post{}, err
	}

	// we set this to be used in service
	post.Project = project

	return post, nil
}

func (s *PostStore) GetPost(ctx context.Context, req request.GetPostRequest) (utils.PaginateResult[entity.Post], error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := s.gormDb.GormDb.WithContext(ctx).
		Model(&entity.Post{}).
		Where("project_id = ?", req.ProjectId).
		Preload("Project", nil).
		Joins("INNER JOIN projects ON projects.id = posts.project_id")

	var searchAllQuery string

	if req.SearchAll != "" {
		searchAllQuery = `
		posts.title ILIKE ?
		OR posts.category ILIKE ?
		OR projects.name ILIKE ?
		OR projects.slug ILIKE ?
		`
	}

	result := utils.ApplyPagination[entity.Post](query, req.PaginationRequest, searchAllQuery)

	if result.Error != nil {
		return utils.PaginateResult[entity.Post]{}, result.Error
	}

	return result, nil
}
