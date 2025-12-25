package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"go.uber.org/zap"

	"github.com/ariefzainuri96/go-logstream/internal/db"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
	"github.com/ariefzainuri96/go-logstream/cmd/api/middleware"
	"gorm.io/gorm"
)

type ProjectStore struct {
	db *db.GormDB
	logger *zap.Logger
}

func (s *ProjectStore) CheckSlug(ctx context.Context, req request.CheckSlugRequest) (bool, error) {
	var slugExists bool

	err := s.db.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.
			Model(&entity.Project{}).
			Select("1"). // return 1 if slug exists (this is signal that row exists)
			Where("slug = ?", req.Slug).
			Limit(1).          // stop query when row found
			Scan(&slugExists). // the destination value is bool, and sql convert value from "1" to true
			Error
	})

	if err != nil {
		return false, err
	}

	return slugExists, nil
}

func (s *ProjectStore) AddProject(ctx context.Context, userId uint, req request.AddProjectRequest) (entity.Project, error) {
	// 1. Retrieve the value (it returns 'any', so you might need to assert it)
    reqID, ok := ctx.Value(middleware.CtxRequestID).(string)
    
    // 2. Safety check: Context values are optional!
    if !ok {
        reqID = "unknown-request" // Fallback if missing
    }

	slugExist, err := s.CheckSlug(ctx, request.CheckSlugRequest{
		Slug: req.Slug,
	})

	if err != nil {
		return entity.Project{}, err
	}

	if slugExist {
		s.logger.Warn("Slug sudah terdaftar", zap.String("RequestId", reqID))
		return entity.Project{}, errors.New("Slug sudah terdaftar")
	}

	project := entity.Project{
		UserId:     userId,
		Name:       req.Name,
		Slug:       req.Slug,
		WebhookUrl: req.WebhookUrl,
	}

	err = s.db.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.
			Create(&project).
			Error
	})

	return project, nil
}

func (s *ProjectStore) GetProject(ctx context.Context, userId uint, req request.PaginationRequest) (utils.PaginateResult[entity.Project], error) {
	var projects []entity.Project

	query := s.db.GormDb.Where("user_id = ?", userId).Find(&projects)

	var searchAllQuery string

	if req.SearchAll != "" {
		searchAllQuery = `
		projects.name ILIKE ?
		OR projects.slug ILIKE ?		
		`
	}

	result := utils.ApplyPagination[entity.Project](query, req, searchAllQuery)

	if result.Error != nil {
		return utils.PaginateResult[entity.Project]{}, result.Error
	}

	return result, nil
}

func (s *ProjectStore) DeleteProject(ctx context.Context, projectId uint) error {
	err := s.db.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.
			Delete(&entity.Project{}, projectId).
			Error
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectStore) UpdateProject(ctx context.Context, projectId uint, req request.AddProjectRequest) (entity.Project, error) {
	project := entity.Project{		
		Name:       req.Name,
		Slug:       req.Slug,
		WebhookUrl: req.WebhookUrl,
	}

	result := s.db.ExecWithTimeoutVal(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.
			Model(&entity.Project{}).
			Where("id = ?", projectId).
			Updates(project)
	})

	if result.Error != nil {
		return entity.Project{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.Project{}, fmt.Errorf("No product found with id %v", projectId)
	}

	if err := s.db.GormDb.First(&project, projectId).Error; err != nil {
        return entity.Project{}, err
    }

	return project, nil
}
