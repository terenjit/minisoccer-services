package services

import (
	"context"
	"field-service/common/util"
	"field-service/constants"
	errFieldSchedule "field-service/constants/error/field_schedule"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FieldScheduleService struct {
	repository repositories.IRepositoryRegistry
}

type IFieldScheduleService interface {
	GetAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) (*util.PaginationResult, error)
	GetAllFieldIdAndDate(context.Context, string, string) ([]dto.FieldScheduleForBookResponse, error)
	GetByUUID(context.Context, string) (*dto.FieldScheduleResponse, error)
	GenerateScheduleForOneMonth(context.Context, *dto.GenerateFieldScheduleOneMonthRequest) error
	Create(context.Context, *dto.FieldScheduleRequest) error
	Update(context.Context, string, *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleResponse, error)
	UpdateStatus(context.Context, *dto.UpdatStatuseFieldScheduleRequest) error
	Delete(context.Context, string) error
}

func NewFieldScheduleService(repository repositories.IRepositoryRegistry) IFieldScheduleService {
	return &FieldScheduleService{repository: repository}
}

func (s *FieldScheduleService) GetAllWithPagination(ctx context.Context, req *dto.FieldScheduleRequestParam) (*util.PaginationResult, error) {
	FieldSchedules, total, err := s.repository.GetFieldSchedule().FindAllWithPagination(ctx, req)
	if err != nil {
		return nil, err
	}

	FieldSchedulesResult := make([]dto.FieldScheduleResponse, 0, len(FieldSchedules))
	for _, FieldSchedule := range FieldSchedules {
		FieldSchedulesResult = append(FieldSchedulesResult, dto.FieldScheduleResponse{
			UUID:         FieldSchedule.UUID,
			FieldName:    FieldSchedule.Field.Name,
			PricePerHour: FieldSchedule.Field.PricePerHour,
			Date:         FieldSchedule.Date.Format("2006-01-02"),
			Status:       FieldSchedule.Status.GetStatusString(),
			Time:         fmt.Sprintf("%s - %s", FieldSchedule.Time.StartTime, FieldSchedule.Time.EndTime),
			CreatedAt:    FieldSchedule.CreatedAt,
			UpdateAt:     FieldSchedule.UpdatdeAt,
		})
	}

	pagination := util.PaginationParam{
		Count: total,
		Page:  req.Page,
		Limit: req.Limit,
		Data:  FieldSchedulesResult,
	}

	response := util.GeneratePagination(pagination)
	return &response, nil
}

func (s *FieldScheduleService) converOneMonthName(inputmonth string) string {
	date, err := time.Parse(time.DateOnly, inputmonth)
	if err != nil {
		return ""
	}

	indonesiaMonth := map[string]string{
		"Jan": "Jan",
		"Feb": "Feb",
		"Mar": "Mar",
		"Apr": "Apr",
		"May": "Mei",
		"Jun": "Jun",
		"Jul": "Jul",
		"Aug": "Agu",
		"Sep": "Sep",
		"Oct": "Okt",
		"Nov": "Nov",
		"Dec": "Des",
	}

	formatedDate := date.Format("02 Jan")
	day := formatedDate[:3]
	month := formatedDate[3:]
	formatedDate = fmt.Sprintf("%s %s", day, indonesiaMonth[month])
	return formatedDate
}

func (s *FieldScheduleService) GetAllFieldIdAndDate(ctx context.Context, uuid string, date string) ([]dto.FieldScheduleForBookResponse, error) {
	FieldSchedules, err := s.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldSchedules, err := s.repository.GetFieldSchedule().FindAllWithFieldIdAndDate(ctx, int(FieldSchedules.ID), date)
	if err != nil {
		return nil, err
	}

	fieldSchedulesResult := make([]dto.FieldScheduleForBookResponse, 0, len(fieldSchedules))
	for _, v := range fieldSchedules {
		priceperHour := float64(v.Field.PricePerHour)
		fieldSchedulesResult = append(fieldSchedulesResult, dto.FieldScheduleForBookResponse{
			UUID:         v.UUID,
			PricePerHour: util.RupiahFormat(&priceperHour),
			Date:         s.converOneMonthName(v.Date.Format(time.DateOnly)),
			Status:       v.Status.GetStatusString(),
			Time:         fmt.Sprintf("%s - %s", v.Time.StartTime, v.Time.EndTime),
		})
	}

	return fieldSchedulesResult, nil

}
func (s *FieldScheduleService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldScheduleResponse, error) {
	FieldSchedule, err := s.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	FieldSchedulesResult := new(dto.FieldScheduleResponse)
	FieldSchedulesResult.UUID = FieldSchedule.UUID
	FieldSchedulesResult.FieldName = FieldSchedule.Field.Name
	FieldSchedulesResult.PricePerHour = FieldSchedule.Field.PricePerHour
	FieldSchedulesResult.Date = s.converOneMonthName(FieldSchedule.Date.Format(time.DateOnly))
	FieldSchedulesResult.Status = FieldSchedule.Status.GetStatusString()
	FieldSchedulesResult.CreatedAt = FieldSchedule.CreatedAt
	FieldSchedulesResult.UpdateAt = FieldSchedule.UpdatdeAt
	FieldSchedulesResult.Time = fmt.Sprintf("%s - %s", FieldSchedule.Time.StartTime, FieldSchedule.Time.EndTime)

	return FieldSchedulesResult, nil

}

func (s *FieldScheduleService) Create(ctx context.Context, req *dto.FieldScheduleRequest) error {
	Field, err := s.repository.GetField().FindByUUID(ctx, req.FieldID)
	if err != nil {
		return err
	}

	fieldSchedules := make([]models.FieldSchedule, 0, len(req.TimeIDs))
	dateParsed, _ := time.Parse(time.DateOnly, req.Date)
	for _, timeId := range req.TimeIDs {
		scheduleTime, err := s.repository.GetTime().FindByUUID(ctx, timeId)
		if err != nil {
			return err
		}

		schedule, err := s.repository.GetFieldSchedule().FindByDateAndTimeId(ctx, req.Date, int(scheduleTime.ID), int(Field.ID))
		if err != nil {
			return err
		}

		if schedule != nil {
			return errFieldSchedule.ErrFieldScheduleExists
		}

		fieldSchedules = append(fieldSchedules, models.FieldSchedule{
			UUID:    uuid.New(),
			FieldID: Field.ID,
			TimeID:  scheduleTime.ID,
			Date:    dateParsed,
			Status:  constants.Available,
		})
	}

	err = s.repository.GetFieldSchedule().Create(ctx, fieldSchedules)
	if err != nil {
		return err
	}

	return nil

}

func (s *FieldScheduleService) GenerateScheduleForOneMonth(ctx context.Context, req *dto.GenerateFieldScheduleOneMonthRequest) error {
	Field, err := s.repository.GetField().FindByUUID(ctx, req.FieldID)
	if err != nil {
		return err
	}

	times, err := s.repository.GetTime().FindAll(ctx)
	if err != nil {
		return err
	}

	numberOfDays := 30
	fieldSchedules := make([]models.FieldSchedule, 0, numberOfDays)
	now := time.Now().Add(time.Duration(1) * 24 * time.Hour)
	for i := 0; i < numberOfDays; i++ {
		currentDate := now.AddDate(0, 0, i)
		for _, v := range times {
			schedule, err := s.repository.GetFieldSchedule().FindByDateAndTimeId(ctx, currentDate.Format(time.DateOnly), int(v.ID), int(Field.ID))
			if err != nil {
				return err
			}

			if schedule != nil {
				return errFieldSchedule.ErrFieldScheduleExists
			}

			fieldSchedules = append(fieldSchedules, models.FieldSchedule{
				UUID:    uuid.New(),
				FieldID: Field.ID,
				TimeID:  v.ID,
				Date:    currentDate,
				Status:  constants.Available,
			})

		}
	}

	err = s.repository.GetFieldSchedule().Create(ctx, fieldSchedules)
	if err != nil {
		return err
	}

	return nil
}

func (s *FieldScheduleService) Update(ctx context.Context, uuid string, req *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleResponse, error) {

	FieldSchedule, err := s.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	scheduleTime, err := s.repository.GetTime().FindByUUID(ctx, req.TimeID)
	if err != nil {
		return nil, err
	}

	isTimeExist, err := s.repository.GetFieldSchedule().FindByDateAndTimeId(ctx, req.Date, int(scheduleTime.ID), int(FieldSchedule.Field.ID))
	if err != nil {
		return nil, err
	}

	if isTimeExist != nil && req.Date != FieldSchedule.Date.Format(time.DateOnly) {
		checkDate, err := s.repository.GetFieldSchedule().FindByDateAndTimeId(ctx, req.Date, int(scheduleTime.ID), int(FieldSchedule.Field.ID))
		if err != nil {
			return nil, err
		}

		if checkDate != nil {
			return nil, errFieldSchedule.ErrFieldScheduleExists
		}
	}

	dateParsed, _ := time.Parse(time.DateOnly, req.Date)
	fieldRes, err := s.repository.GetFieldSchedule().Update(ctx, uuid, &models.FieldSchedule{
		Date:   dateParsed,
		TimeID: scheduleTime.ID,
	})
	if err != nil {
		return nil, err
	}

	response := dto.FieldScheduleResponse{
		UUID:         fieldRes.UUID,
		FieldName:    fieldRes.Field.Name,
		Date:         fieldRes.Date.Format(time.DateOnly),
		PricePerHour: fieldRes.Field.PricePerHour,
		Status:       FieldSchedule.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", scheduleTime.StartTime, scheduleTime.EndTime),
		CreatedAt:    fieldRes.CreatedAt,
		UpdateAt:     fieldRes.UpdatdeAt,
	}

	return &response, nil
}

func (s *FieldScheduleService) UpdateStatus(ctx context.Context, req *dto.UpdatStatuseFieldScheduleRequest) error {

	for _, item := range req.FiledSchedulesIDs {
		_, err := s.repository.GetFieldSchedule().FindByUUID(ctx, item)
		if err != nil {
			return err
		}

		err = s.repository.GetFieldSchedule().UpdateStatus(ctx, constants.Booked, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *FieldScheduleService) Delete(ctx context.Context, uuid string) error {

	_, err := s.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.repository.GetFieldSchedule().Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}
