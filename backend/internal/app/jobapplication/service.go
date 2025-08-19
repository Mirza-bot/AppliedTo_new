package jobapplication

import (
	"appliedTo/internal/platform/validate"
	"context"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CREATE
func (s *Service) Create(ctx context.Context, in JobApplicationCreateDto) (JobApplicationPublicDto, error) {
	if err := validate.Required(
		validate.Field{Name: "title",             Value: in.BaseJobApplicationDto.Title},
		validate.Field{Name: "employment type",   Value: in.BaseJobApplicationDto.Employment.Type},
		validate.Field{Name: "work location",     Value: in.BaseJobApplicationDto.Employment.WorkLocation},
	); err != nil {
		return JobApplicationPublicDto{}, err
	}

	m := CreateModel(in)
	if err := s.db.WithContext(ctx).Create(&m).Error; err != nil {
		return JobApplicationPublicDto{}, err
	}
	return MapModelToPublicDto(m), nil
}

// READ
func (s *Service) GetByID(ctx context.Context, id uint) (JobApplicationPublicDto, error) {
	var m JobApplication
	if err := s.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return JobApplicationPublicDto{}, err
	}
	return MapModelToPublicDto(m), nil
}

// UPDATE (full replace)
func (s *Service) Update(ctx context.Context, id uint, in JobApplicationCreateDto) (JobApplicationPublicDto, error) {
	// ensure record exists
	var existing JobApplication
	if err := s.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		return JobApplicationPublicDto{}, err
	}

	if err := validate.Required(
		validate.Field{Name: "title",             Value: in.BaseJobApplicationDto.Title},
		validate.Field{Name: "employment type",   Value: in.BaseJobApplicationDto.Employment.Type},
		validate.Field{Name: "work location",     Value: in.BaseJobApplicationDto.Employment.WorkLocation},
	); err != nil {
		return JobApplicationPublicDto{}, err
	}

	updated := CreateModel(in)       // build a fresh model from DTO
	// apply onto the existing row (existing has the primary key set)
	if err := s.db.WithContext(ctx).Model(&existing).Updates(updated).Error; err != nil {
		return JobApplicationPublicDto{}, err
	}

	// reload if you need updated associations; otherwise map 'existing'
	if err := s.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		return JobApplicationPublicDto{}, err
	}
	return MapModelToPublicDto(existing), nil
}

// PATCH (partial update)
func (s *Service) Patch(ctx context.Context, id uint, patch JobApplicationPatchDto) (JobApplicationPublicDto, error) {
	var m JobApplication
	if err := s.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return JobApplicationPublicDto{}, err
	}

	PatchModel(&m, patch)

	if err := s.db.WithContext(ctx).Save(&m).Error; err != nil {
		return JobApplicationPublicDto{}, err
	}

	return MapModelToPublicDto(m), nil
}

// DELETE
func (s *Service) Delete(ctx context.Context, id uint) error {
	tx := s.db.WithContext(ctx).Delete(&JobApplication{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
