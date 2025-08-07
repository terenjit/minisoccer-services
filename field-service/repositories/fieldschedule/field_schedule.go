package repositories

import (
	"context"
	"errors"
	errWrap "field-service/common/error"
	"field-service/constants"
	errConstant "field-service/constants/error"
	errFieldSchedule "field-service/constants/error/field_schedule"
	"field-service/domain/dto"
	"field-service/domain/models"
	"fmt"

	"gorm.io/gorm"
)

type FieldScheduleRepository struct {
	db *gorm.DB
}

type IFieldScheduleRepository interface {
	FindAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error)
	FindAllWithFieldIdAndDate(context.Context, int, string) ([]models.FieldSchedule, error)
	FindByUUID(context.Context, string) (*models.FieldSchedule, error)
	FindByDateAndTimeId(context.Context, string, int, int) (*models.FieldSchedule, error)
	Create(context.Context, []models.FieldSchedule) error
	Update(context.Context, string, *models.FieldSchedule) (*models.FieldSchedule, error)
	UpdateStatus(context.Context, constants.FieldScheduleStatus, string) error
	Delete(context.Context, string) error
}

func NewFieldScheduleRepository(db *gorm.DB) IFieldScheduleRepository {
	return &FieldScheduleRepository{db: db}
}

func (f *FieldScheduleRepository) FindAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error) {
	var (
		fields []models.FieldSchedule
		sort   string
		total  int64
	)

	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *&param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit

	err := f.db.WithContext(ctx).Preload("Field").Preload("Time").Limit(limit).Offset(offset).Order(sort).Find(&fields).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	err = f.db.WithContext(ctx).Model(&fields).Count(&total).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return fields, total, nil

}

func (f *FieldScheduleRepository) FindAllWithFieldIdAndDate(ctx context.Context, id int, time string) ([]models.FieldSchedule, error) {
	var (
		fields []models.FieldSchedule
	)

	err := f.db.WithContext(ctx).Preload("Field").Preload("Time").Where("field_id = ?", id).Where("date = ?", time).
		Joins("LEFT JOIN times on field_schedules.time_id = times.id").Order("times.start_time asc").
		Find(&fields).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return fields, nil

}

func (f *FieldScheduleRepository) FindByUUID(ctx context.Context, uuid string) (*models.FieldSchedule, error) {
	var (
		field *models.FieldSchedule
	)

	err := f.db.WithContext(ctx).Preload("Field").Preload("Time").Where("uuid = ?", uuid).First(&field).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return field, nil
}

func (f *FieldScheduleRepository) FindByDateAndTimeId(ctx context.Context, date string, timeId int, fieldId int) (*models.FieldSchedule, error) {
	var (
		field *models.FieldSchedule
	)

	err := f.db.WithContext(ctx).Where("date = ?", date).Where("time_id = ?", timeId).Where("field_id = ?", fieldId).First(&field).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return field, nil
}

func (f *FieldScheduleRepository) Create(ctx context.Context, req []models.FieldSchedule) error {
	err := f.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil
}

func (f *FieldScheduleRepository) Update(ctx context.Context, uuid string, req *models.FieldSchedule) (*models.FieldSchedule, error) {
	fieldSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldSchedule.Date = req.Date

	err = f.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return fieldSchedule, nil

}

func (f *FieldScheduleRepository) UpdateStatus(ctx context.Context, status constants.FieldScheduleStatus, uuid string) error {
	fieldSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	fieldSchedule.Status = status

	err = f.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil

}

func (f *FieldScheduleRepository) Delete(ctx context.Context, uuid string) error {
	var (
		fieldSchedule *models.FieldSchedule
	)

	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&fieldSchedule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil

}
