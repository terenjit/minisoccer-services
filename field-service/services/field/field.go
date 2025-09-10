package services

import (
	"bytes"
	"context"
	"field-service/common/gcs"
	"field-service/common/util"
	errConstant "field-service/constants/error"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type FieldService struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
}

type IfieldService interface {
	GetAllWithPagination(context.Context, *dto.FieldRequestParam) (*util.PaginationResult, error)
	GetAllWithoutPagination(context.Context) ([]dto.FieldResponse, error)
	GetByUUID(context.Context, string) (*dto.FieldResponse, error)
	Create(context.Context, *dto.FieldRequest) (*dto.FieldResponse, error)
	Update(context.Context, string, *dto.UpdateFieldRequest) (*dto.FieldResponse, error)
	Delete(context.Context, string) error
}

func NewFieldService(repository repositories.IRepositoryRegistry, gcs gcs.IGCSClient) IfieldService {
	return &FieldService{repository: repository, gcs: gcs}
}

func (s *FieldService) GetAllWithPagination(ctx context.Context, req *dto.FieldRequestParam) (*util.PaginationResult, error) {
	fields, total, err := s.repository.GetField().FindAllWithPagination(ctx, req)
	if err != nil {
		return nil, err
	}

	fieldsResult := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldsResult = append(fieldsResult, dto.FieldResponse{
			UUID:         field.UUID,
			Code:         field.Code,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Images,
			CreatedAt:    field.CreatedAt,
			UpdateAt:     field.UpdatedAt,
		})
	}

	pagination := util.PaginationParam{
		Count: total,
		Page:  req.Page,
		Limit: req.Limit,
		Data:  fieldsResult,
	}

	response := util.GeneratePagination(pagination)
	return &response, nil
}

func (s *FieldService) GetAllWithoutPagination(ctx context.Context) ([]dto.FieldResponse, error) {
	fields, err := s.repository.GetField().FindAllWithoutPagination(ctx)
	if err != nil {
		return nil, err
	}

	fieldsResult := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldsResult = append(fieldsResult, dto.FieldResponse{
			UUID:         field.UUID,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Images,
		})
	}

	return fieldsResult, nil

}
func (s *FieldService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldResponse, error) {
	field, err := s.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldsResult := new(dto.FieldResponse)
	fieldsResult.UUID = field.UUID
	fieldsResult.Code = field.Code
	fieldsResult.Name = field.Name
	fieldsResult.PricePerHour = field.PricePerHour
	fieldsResult.Images = field.Images
	fieldsResult.CreatedAt = field.CreatedAt
	fieldsResult.UpdateAt = field.UpdatedAt

	return fieldsResult, nil

}

func (s *FieldService) validateUpload(images []multipart.FileHeader) error {
	if len(images) <= 0 {
		return errConstant.ErrInvalidUploadFile
	}

	for _, image := range images {
		if image.Size > 5*1024*1024 {
			return errConstant.ErrSizetooBig
		}
	}
	return nil
}

func (s *FieldService) processAndUpload(ctx context.Context, image multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("images/%s-%s-%s", time.Now().Format("20060102150405"), image.Filename, path.Ext(image.Filename))
	url, err := s.gcs.UploadFile(ctx, fileName, buffer.Bytes())
	if err != nil {
		return "", nil
	}

	return url, nil
}

func (s *FieldService) uploadImages(ctx context.Context, images []multipart.FileHeader) ([]string, error) {
	err := s.validateUpload(images)
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0, len(images))
	for _, v := range images {
		url, err := s.processAndUpload(ctx, v)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (s *FieldService) Create(ctx context.Context, req *dto.FieldRequest) (*dto.FieldResponse, error) {
	imageUrl, err := s.uploadImages(ctx, req.Images)
	if err != nil {
		return nil, err
	}

	field, err := s.repository.GetField().Create(ctx, &models.Field{
		Code:         req.Code,
		Name:         req.Name,
		PricePerHour: req.PricePerHour,
		Images:       pq.StringArray(imageUrl),
	})
	if err != nil {
		logrus.Errorf("error create field: %v", err)
		return nil, err
	}

	response := &dto.FieldResponse{
		UUID:         field.UUID,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: field.PricePerHour,
		Images:       field.Images,
		CreatedAt:    field.CreatedAt,
		UpdateAt:     field.UpdatedAt,
	}

	return response, nil

}
func (s *FieldService) Update(ctx context.Context, uuidParam string, req *dto.UpdateFieldRequest) (*dto.FieldResponse, error) {
	var image []string

	field, err := s.repository.GetField().FindByUUID(ctx, uuidParam)
	if err != nil {
		return nil, err
	}

	if req.Images == nil {
		image = field.Images
	} else {
		imageUrl, err := s.uploadImages(ctx, req.Images)
		if err != nil {
			return nil, err
		}
		image = append(image, imageUrl...)
	}

	field, err = s.repository.GetField().Update(ctx, uuidParam, &models.Field{
		Code:         req.Code,
		Name:         req.Name,
		PricePerHour: req.PricePerHour,
		Images:       image,
	})
	if err != nil {
		return nil, err
	}
	uuidParsed, _ := uuid.Parse(uuidParam)
	response := &dto.FieldResponse{
		UUID:         uuidParsed,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: field.PricePerHour,
		Images:       field.Images,
		CreatedAt:    field.CreatedAt,
		UpdateAt:     field.UpdatedAt,
	}

	return response, nil
}
func (s *FieldService) Delete(ctx context.Context, uuid string) error {

	_, err := s.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.repository.GetField().Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}
