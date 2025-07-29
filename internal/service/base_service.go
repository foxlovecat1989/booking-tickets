package service

import (
	"tickets/internal/repository"
)

// BaseService provides common service functionality
type BaseService struct {
	baseRepo *repository.BaseRepository
}

// NewBaseService creates a new base service
func NewBaseService(baseRepo *repository.BaseRepository) *BaseService {
	return &BaseService{
		baseRepo: baseRepo,
	}
}

// GetBaseRepository returns the base repository instance
func (s *BaseService) GetBaseRepository() *repository.BaseRepository {
	return s.baseRepo
}
