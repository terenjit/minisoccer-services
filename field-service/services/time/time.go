package services

import (
	"context"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
)

type TimeService struct {
	repository repositories.IRepositoryRegistry
}

type ITimeService interface {
	GetAll(context.Context) ([]dto.TimeResponse, error)
	GetByUUID(context.Context, string) (*dto.TimeResponse, error)
	Create(context.Context, *dto.TimeRequest) (*dto.TimeResponse, error)
}

func NewTimeService(repository repositories.IRepositoryRegistry) ITimeService {
	return &TimeService{repository: repository}
}

func (s *TimeService) GetAll(ctx context.Context) ([]dto.TimeResponse, error) {

	times, err := s.repository.GetTime().FindAll(ctx)
	if err != nil {
		return nil, err
	}

	timeResults := make([]dto.TimeResponse, 0, len(times))
	for _, time := range times {
		timeResults = append(timeResults, dto.TimeResponse{
			UUID:      time.UUID,
			StartTime: time.StartTime,
			EndTime:   time.EndTime,
			CreatedAt: time.CreatedAt,
			UpdateAt:  time.UpdatedAt,
		})
	}

	return timeResults, nil
}

func (s *TimeService) GetByUUID(ctx context.Context, id string) (*dto.TimeResponse, error) {

	time, err := s.repository.GetTime().FindByUUID(ctx, id)
	if err != nil {
		return nil, err
	}

	timeResult := dto.TimeResponse{
		UUID:      time.UUID,
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
		CreatedAt: time.CreatedAt,
		UpdateAt:  time.UpdatedAt,
	}

	return &timeResult, nil
}

func (s *TimeService) Create(ctx context.Context, req *dto.TimeRequest) (*dto.TimeResponse, error) {

	time := &models.Time{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	time, err := s.repository.GetTime().Create(ctx, time)
	if err != nil {
		return nil, err
	}

	timeResult := dto.TimeResponse{
		UUID:      time.UUID,
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
		CreatedAt: time.CreatedAt,
		UpdateAt:  time.UpdatedAt,
	}

	return &timeResult, nil
}
