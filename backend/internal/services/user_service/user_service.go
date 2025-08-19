package userservice

import (
	userdtos "appliedTo/dtos/user_dtos"
	"appliedTo/internal/validate"
	mappers "appliedTo/mappers/user_mappers"
	"appliedTo/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrEmailInUse   = errors.New("email already in use")
	ErrInvalidEmail = errors.New("invalid email")
)

type UserService struct {
	db     *gorm.DB
	hasher interface {
		Hash(string) (string, error)
	}
}

func NewUserService(db *gorm.DB, hasher interface {
	Hash(string) (string, error)
}) *UserService {
	return &UserService{db: db, hasher: hasher}
}

// -------- CREATE --------

func (s *UserService) Create(ctx context.Context, dto userdtos.UserCreateDto) (userdtos.UserPublicDto, error) {
	if err := validate.Required(
		validate.Field{Name: "firstname", Value: dto.FirstName},
		validate.Field{Name: "lastname", Value: dto.LastName},
		validate.Field{Name: "email", Value: dto.Email},
		validate.Field{Name: "password", Value: dto.Password},
	); err != nil {
		return userdtos.UserPublicDto{}, err
	}

	normalizedEmail, err := validate.NormalizeAndValidateEmail(dto.Email)
	if err != nil {
		return userdtos.UserPublicDto{}, ErrInvalidEmail
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", normalizedEmail).Count(&count).Error; err != nil {
		return userdtos.UserPublicDto{}, ErrEmailInUse
	}

	hash, err := s.hasher.Hash(dto.Password)
	if err != nil {
		return userdtos.UserPublicDto{}, err
	}

	user := mappers.CreateModel(dto)
	user.Email = normalizedEmail
	user.Password = hash

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return userdtos.UserPublicDto{}, ErrEmailInUse
		}
		return userdtos.UserPublicDto{}, err
	}

	return mappers.MapModelToPublicDto(user), nil
}

// -------- READ --------

func (s *UserService) GetByID(ctx context.Context, id uint) (userdtos.UserPublicDto, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return userdtos.UserPublicDto{}, err
	}
	return mappers.MapModelToPublicDto(user), nil
}

// -------- UPDATE --------

func (s *UserService) Update(ctx context.Context, id uint, dto userdtos.UserCreateDto) (userdtos.UserPublicDto, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return userdtos.UserPublicDto{}, err
	}

	if err := validate.Required(
		validate.Field{Name: "firstname", Value: dto.FirstName},
		validate.Field{Name: "lastname", Value: dto.LastName},
		validate.Field{Name: "email", Value: dto.Email},
		validate.Field{Name: "password", Value: dto.Password},
	); err != nil {
		return userdtos.UserPublicDto{}, err
	}

	normalizedEmail, err := validate.NormalizeAndValidateEmail(dto.Email)
	if err != nil {
		return userdtos.UserPublicDto{}, ErrInvalidEmail
	}

	taken, err := s.emailTaken(ctx, normalizedEmail, user.ID)
	if err != nil {
		return userdtos.UserPublicDto{}, err
	}
	if taken {
		return userdtos.UserPublicDto{}, ErrEmailInUse
	}

	hash, err := s.hasher.Hash(dto.Password)
	if err != nil {
		return userdtos.UserPublicDto{}, err
	}

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	user.Email = normalizedEmail
	user.Password = hash

	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return userdtos.UserPublicDto{}, err
	}

	return mappers.MapModelToPublicDto(user), nil
}

// -------- PATCH --------

func (s *UserService) Patch(ctx context.Context, id uint, dto userdtos.UserPatchDto) (userdtos.UserPublicDto, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return userdtos.UserPublicDto{}, err
	}

	if dto.Email != nil {
		norm, err := validate.NormalizeAndValidateEmail(*dto.Email)
		if err != nil {
			return userdtos.UserPublicDto{}, ErrInvalidEmail
		}
		if norm != user.Email {
			taken, err := s.emailTaken(ctx, norm, user.ID)
			if err != nil {
				return userdtos.UserPublicDto{}, err
			}
			if taken {
				return userdtos.UserPublicDto{}, ErrEmailInUse
			}
			user.Email = norm
		}
	}

	mappers.PatchModel(&user, dto)

	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return userdtos.UserPublicDto{}, err
	}

	return mappers.MapModelToPublicDto(user), nil
}

// -------- DELETE --------

func (s *UserService) Delete(ctx context.Context, id uint) error {
	tx := s.db.WithContext(ctx).Delete(&models.User{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	// GORM doesn't error when nothing is deleted; turn that into a 404 upstream if desired
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// -------- helpers --------

func (s *UserService) emailTaken(ctx context.Context, normalizedEmail string, excludeID uint) (bool, error) {
	var count int64
	q := s.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", normalizedEmail)
	if excludeID != 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
