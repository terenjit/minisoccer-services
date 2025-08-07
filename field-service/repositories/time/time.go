package repositories

import (
	"context"
	"errors"
	errWrap "field-service/common/error"
	errConstant "field-service/constants/error"
	errTime "field-service/constants/error/time"
	"field-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeRepository struct {
	db *gorm.DB
}

type ITimeRepository interface {
	FindAll(context.Context) ([]models.Time, error)
	FindByUUID(context.Context, string) ([]models.Time, error)
	FindID(context.Context, int) ([]models.Time, error)
	Create(context.Context, *models.Time) (*models.Time, error)
}

func NewTimeRepository(db *gorm.DB) ITimeRepository {
	return &TimeRepository{db: db}
}

func (t *TimeRepository) FindAll(ctx context.Context) ([]models.Time, error) {
	var times []models.Time
	err := t.db.WithContext(ctx).Find(&times).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errTime.ErrTimeNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return times, err
}

func (t *TimeRepository) FindByUUID(ctx context.Context, uuid string) ([]models.Time, error) {
	var times []models.Time
	err := t.db.WithContext(ctx).Where("uuid = ?", uuid).Find(&times).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errTime.ErrTimeNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return times, err
}

func (t *TimeRepository) FindID(ctx context.Context, id int) ([]models.Time, error) {
	var times []models.Time
	err := t.db.WithContext(ctx).Where("id = ?", id).Find(&times).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errTime.ErrTimeNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return times, err
}

func (t *TimeRepository) Create(ctx context.Context, req *models.Time) (*models.Time, error) {
	time := models.Time{
		UUID: uuid.New(),
	}

	err := t.db.WithContext(ctx).Create(&time).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errTime.ErrTimeNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &time, nil
}
