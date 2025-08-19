package user

import (
	"appliedTo/internal/platform/validate"
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrEmailInUse   = errors.New("email already in use")
	ErrInvalidEmail = errors.New("invalid email")
)

type Service struct {
	db     *gorm.DB
	hasher interface {
		Hash(string) (string, error)
	}
}

func NewService(db *gorm.DB, hasher interface {
	Hash(string) (string, error)
}) *Service {
	return &Service{db: db, hasher: hasher}
}

// -------- CREATE --------

func (s *Service) Create(ctx context.Context, dto UserCreateDto) (UserPublicDto, error) {
	if err := validate.Required(
		validate.Field{Name: "firstname", Value: dto.FirstName},
		validate.Field{Name: "lastname", Value: dto.LastName},
		validate.Field{Name: "email", Value: dto.Email},
		validate.Field{Name: "password", Value: dto.Password},
	); err != nil {
		return UserPublicDto{}, err
	}

	normalizedEmail, err := validate.NormalizeAndValidateEmail(dto.Email)
	if err != nil {
		return UserPublicDto{}, ErrInvalidEmail
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&User{}).Where("email = ?", normalizedEmail).Count(&count).Error; err != nil {
		return UserPublicDto{}, ErrEmailInUse
	}

	hash, err := s.hasher.Hash(dto.Password)
	if err != nil {
		return UserPublicDto{}, err
	}

	user := CreateModel(dto)
	user.Email = normalizedEmail
	user.Password = hash

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return UserPublicDto{}, ErrEmailInUse
		}
		return UserPublicDto{}, err
	}

	return MapModelToPublicDto(user), nil
}

// -------- READ --------

func (s *Service) GetByID(ctx context.Context, id uint) (UserPublicDto, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return UserPublicDto{}, err
	}
	return MapModelToPublicDto(user), nil
}

// -------- UPDATE --------

func (s *Service) Update(ctx context.Context, id uint, dto UserCreateDto) (UserPublicDto, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return UserPublicDto{}, err
	}

	if err := validate.Required(
		validate.Field{Name: "firstname", Value: dto.FirstName},
		validate.Field{Name: "lastname", Value: dto.LastName},
		validate.Field{Name: "email", Value: dto.Email},
		validate.Field{Name: "password", Value: dto.Password},
	); err != nil {
		return UserPublicDto{}, err
	}

	normalizedEmail, err := validate.NormalizeAndValidateEmail(dto.Email)
	if err != nil {
		return UserPublicDto{}, ErrInvalidEmail
	}

	taken, err := s.emailTaken(ctx, normalizedEmail, user.ID)
	if err != nil {
		return UserPublicDto{}, err
	}
	if taken {
		return UserPublicDto{}, ErrEmailInUse
	}

	hash, err := s.hasher.Hash(dto.Password)
	if err != nil {
		return UserPublicDto{}, err
	}

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	user.Email = normalizedEmail
	user.Password = hash

	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return UserPublicDto{}, err
	}

	return MapModelToPublicDto(user), nil
}

// -------- PATCH --------

func (s *Service) Patch(ctx context.Context, id uint, dto UserPatchDto) (UserPublicDto, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return UserPublicDto{}, err
	}

	if dto.Email != nil {
		norm, err := validate.NormalizeAndValidateEmail(*dto.Email)
		if err != nil {
			return UserPublicDto{}, ErrInvalidEmail
		}
		if norm != user.Email {
			taken, err := s.emailTaken(ctx, norm, user.ID)
			if err != nil {
				return UserPublicDto{}, err
			}
			if taken {
				return UserPublicDto{}, ErrEmailInUse
			}
			user.Email = norm
		}
	}

	PatchModel(&user, dto)

	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return UserPublicDto{}, err
	}

	return MapModelToPublicDto(user), nil
}

// -------- DELETE --------

func (s *Service) Delete(ctx context.Context, id uint) error {
	tx := s.db.WithContext(ctx).Delete(&User{}, id)
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

func (s *Service) emailTaken(ctx context.Context, normalizedEmail string, excludeID uint) (bool, error) {
	var count int64
	q := s.db.WithContext(ctx).Model(&User{}).Where("email = ?", normalizedEmail)
	if excludeID != 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
